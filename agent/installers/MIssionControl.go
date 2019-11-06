package installers

import (
	agentUtils "agent/utils"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
	"github.com/thoas/go-funk"
	"strings"
)

func InstallMissionControl(name string, setupInfo structs.SetupInfo) {
	logger.Always("Install Mission Control")

	commonUtils.HelmInstall(name,"jfrog/mission-control" , setupInfo.Services.Versions.Jfmc , []string {
		"--set mongodb.db.adminPassword=zooloo " ,
		"--set mongodb.db.insightPassword=zooloo " ,
		"--set mongodb.db.mcPassword=zooloo " ,
		"--set postgresql.db.jfisPassword=zooloo " ,
		"--set postgresql.db.jfscPassword=zooloo " ,
		"--set postgresql.db.jfexPassword=zooloo " ,
		"--set postgresql.db.jfmcPassword=zooloo " ,
		"--set postgresql.persistence.size=10Gi " ,
		"--set elasticsearch.persistence.size=10Gi " ,
		"--set missionControl.persistence.size=10Gi " ,
		"--set postgresql.postgresPassword=zooloo " } ,"")

	agentUtils.WaitForService(name, 80, true)
}

func SetupMissionControl(setupInfo structs.SetupInfo) {
	logger.Always("Setup Mission Control")
	setMissionControlBaseUrl()
	loadLicenseBuckets(setupInfo)
	addJfmcSites(setupInfo)
}

func AddRepositories(servers []structs.Artifactory) {
	logger.Always("Add artifactory repositories")

	for _, server := range servers {
		for _ ,repo := range server.Repos {
			addRepository(server.Name ,repo)
		}
	}
}

func SetupReplications(servers []structs.Artifactory) {
	logger.Always("Setup Replications")

	for _, server := range servers {
		for _, replication := range server.Replications {
			AddReplication("star-replication" ,replication.Source ,server.Name ,replication.Repo)
		}
	}
}

func AddFederationRules(servers []structs.Artifactory) {
	logger.Always("Setup Federation rules")

	ip := kubernetes.GetServiceIp("mission-control")

	for _, server := range servers {
		var targets []structs.ArtifactoryHost

		funk.ForEach(servers, func(otherServer structs.Artifactory) {
			if server.Name != otherServer.Name {
				targets = append(targets, structs.ArtifactoryHost{Name: otherServer.Name})
			}
		})

		commonUtils.InvokeRequest("http://"+ip+"/api/v3/services/access/"+server.Name+"/federation", "PUT", structs.AddAccessFederationPayload{
			Entities: []string{"USERS", "GROUPS", "PERMISSIONS", "TOKENS"},
			Targets:  targets,
		})
	}
}

func SetupCircleOfTrust(servers []structs.Artifactory) {
	logger.Always("Setup Circle of trust between all nodes")

	artifactoryHome, _, _ := kubernetes.ExecOnPod("printenv ARTIFACTORY_HOME", servers[0].Name + "-artifactory-0",
		"artifactory", "default", nil)

	artifactoryHome = strings.TrimSuffix(artifactoryHome, "\n")

	for _, server := range servers {
		for _, anotherServer := range servers {
			crt, _, _ := kubernetes.DownloadFileFromPod(artifactoryHome+"/access/etc/keys/root.crt", server.Name+"-artifactory-0", "artifactory", "default")
			_, _, _ = kubernetes.UploadContentToPod([]byte(crt), artifactoryHome + "/access/etc/keys/trusted/" +server.Name+ ".crt", anotherServer.Name+ "-artifactory-0", "artifactory", "default")
		}
	}
}

func setMissionControlBaseUrl() {
	logger.Always("Setup Mission Control base url")

	ip := kubernetes.GetServiceIp("mission-control")

	data := structs.BaseUrlPayload{BaseURL: "http://" + ip}
	commonUtils.InvokeRequest("http://"+ip+"/api/v3/settings/base_url", "PUT", data)
}

func loadLicenseBuckets(setupInfo structs.SetupInfo) {
	logger.Always("Load Mission Control buckets")

	ip := kubernetes.GetServiceIp("mission-control")

	artifactoryParams := map[string]string{
		"bucket_name": "enterprise-licenses",
		"bucket_key":  setupInfo.ArtifactoryLicense,
	}

	err := commonUtils.UploadFile("http://"+ip+"/api/v3/buckets", artifactoryParams, "bucket_file",
		commonUtils.GetResourceLocation("buckets/artifactory.json" ,commonUtils.GetAgentResources()))

	if err != nil {
		logger.Critical("unable to upload license")
	}

	edgeParams := map[string]string{
		"bucket_name": "edge-licenses",
		"bucket_key":  setupInfo.EdgeLicense,
	}

	err = commonUtils.UploadFile("http://"+ip+"/api/v3/buckets", edgeParams, "bucket_file",
		commonUtils.GetResourceLocation("buckets/edge.json" ,commonUtils.GetAgentResources()))

	if err != nil {
		logger.Critical("unable to upload license")
	}
}

func addJfmcSites(setupInfo structs.SetupInfo) {
	logger.Always("Load Mission Control sites")

	ip := kubernetes.GetServiceIp("mission-control")

	for _, site := range setupInfo.Sites {
		commonUtils.InvokeRequest("http://"+ip+"/api/v3/sites", "POST", site)
	}
}

func addRepository(artiName string, repo structs.Repo) {
	artifactoryAddress := kubernetes.GetServiceIp(artiName + artifactoryNginxService)

	if repo.Local {
		commonUtils.InvokeRequestWithContentType("http://"+artifactoryAddress+"/artifactory/api/repositories/" + repo.Name + "-local", "PUT", structs.AddRepoPayload{
			Type: "local",
			PackageType: repo.PackageType,
		} , "application/vnd.org.jfrog.artifactory.repositories.LocalRepositoryConfiguration+json")
	}

	if repo.Remote {
		commonUtils.InvokeRequestWithContentType("http://"+artifactoryAddress+"/artifactory/api/repositories/" + repo.Name + "-remote", "PUT", structs.AddRepoPayload{
			Type: "remote",
			PackageType: repo.PackageType,
			URL: repo.Url,
			ExternalDependencyEnabled: true,
		} ,"application/vnd.org.jfrog.artifactory.repositories.RemoteRepositoryConfiguration+json")
	}

	var repos []string

	if repo.Local {repos = append(repos ,repo.Name + "-local")}
	if repo.Remote {repos = append(repos ,repo.Name + "-remote")}

	commonUtils.InvokeRequestWithContentType("http://"+artifactoryAddress+"/artifactory/api/repositories/" + repo.Name, "PUT", structs.AddRepoPayload{
		Type: "virtual",
		PackageType: repo.PackageType,
		Repositories: repos,
		DefaultDeploymentRepo: repo.Name + "-local",
		ExternalDependencyEnabled: true,
	} ,"application/vnd.org.jfrog.artifactory.repositories.VirtualRepositoryConfiguration+json")

}

func AddReplication(replicationType string ,sourceArti string, targetArti string, repoName string) {
	// todo implement
}

func setMissionControlLoggers() string {
	return " --set missionControl.loggers[0]=jfmc-server.log " +
		" --set missionControl.loggers[1]=audit.log" +
		" --set missionControl.loggers[2]=http.log "
}
