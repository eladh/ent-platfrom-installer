package installers

import (
	"common/cloud/gcp"
	"common/kubernetes"
	"common/structs"
	"common/utils"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const defaultToken string = "apiToken f562d4fa-cf5f-4d58-8b9b-91ef949c024f"

func InstallPipeline(setupInfo structs.SetupInfo, artiAddress string) {
	if setupInfo.Pipeline.GitRepo == "" {
		return
	}

	prefix := setupInfo.Cluster.Name + setupInfo.Cluster.Domain

	wwwServiceIp := gcp.CreateExternalIPAddress(prefix+"-"+"www", setupInfo)
	apiServiceIp := gcp.CreateExternalIPAddress(prefix+"-"+"api", setupInfo)

	InstallPipelineChart(setupInfo, apiServiceIp, wwwServiceIp, artiAddress)
	kubernetes.WaitForPod(setupInfo.Cluster.Name+"-pipelines-k8s-node-0", 2)

	id := AddGitHubPipelineIntegration(apiServiceIp, "git", setupInfo.Pipeline.GitToken)

	AddPipelineSource(apiServiceIp, id, setupInfo.Pipeline.GitRepo, setupInfo.Pipeline.GitSource)
	AddArtifactoryIntegration(apiServiceIp, artiAddress, "art")

	distributionIp := kubernetes.GetServiceIp("distribution")
	AddDistributionIntegration(apiServiceIp, distributionIp, "dist")

	certificate, _ := ioutil.ReadFile(setupInfo.TempDir + "ca.crt")
	dockerRegistryFqdn := "docker." + GetAuthServerFullName(GetAuthProvider(setupInfo).Name, setupInfo)
	updateDockerCertOnPod(certificate, setupInfo.Cluster.Name+"-pipelines-k8s-node-0", dockerRegistryFqdn)
}

func updateDockerCertOnPod(certificate []byte, podName string, dockerRegistryFqdn string) {
	targetPath := "/etc/docker/certs.d/" + dockerRegistryFqdn
	_, _, _ = kubernetes.ExecOnPod("mkdir -p "+targetPath, podName, "dind", "default", nil)
	_, _, _ = kubernetes.UploadContentToPod(certificate, targetPath+"/ca.crt", podName, "dind", "default")

}

func AddArtifactoryIntegration(pipelineServer string, artiServer string, name string) {
	utils.InvokeRequestWithToken("http://"+pipelineServer+":30000/projectIntegrations", "POST",
		structs.PipelineIntegration{
			Name:                  name,
			ProjectID:             1,
			MasterIntegrationID:   98,
			MasterIntegrationName: "artifactory",
			MasterIntegrationType: "generic",
			FormJSONValues: []structs.LabelValuePair{
				{Label: "url", Value: artiServer},
				{Label: "apikey", Value: "password"},
				{Label: "user", Value: "admin"},
			},
		}, defaultToken)
}

func AddSshIntegration(pipelineServer string, name string, privateKey string, publicKey string) {
	utils.InvokeRequestWithToken("http://"+pipelineServer+":30000/projectIntegrations", "POST",
		structs.PipelineIntegration{
			Name:                  name,
			ProjectID:             1,
			MasterIntegrationID:   71,
			MasterIntegrationName: "sshKey",
			MasterIntegrationType: "generic",
			FormJSONValues: []structs.LabelValuePair{
				{Label: "privateKey", Value: privateKey},
				{Label: "publicKey", Value: publicKey},
			},
		}, defaultToken)
}

func AddDistributionIntegration(pipelineServer string, distAddress string, name string) {
	utils.InvokeRequestWithToken("http://"+pipelineServer+":30000/projectIntegrations", "POST",
		structs.PipelineIntegration{
			Name:                  name,
			ProjectID:             1,
			MasterIntegrationID:   100,
			MasterIntegrationName: "distribution",
			MasterIntegrationType: "generic",
			FormJSONValues: []structs.LabelValuePair{
				{Label: "url", Value: distAddress},
				{Label: "apikey", Value: "password"},
				{Label: "user", Value: "admin"},
			},
		}, defaultToken)
}

func AddGitHubPipelineIntegration(pipelineServer string, name string, token string) int {
	result := utils.InvokeRequestWithToken("http://"+pipelineServer+":30000/projectIntegrations", "POST",
		structs.PipelineIntegration{
			Name:                  name,
			ProjectID:             1,
			MasterIntegrationID:   20,
			MasterIntegrationName: "github",
			MasterIntegrationType: "scm",
			FormJSONValues: []structs.LabelValuePair{
				{Label: "url", Value: "https://api.github.com"},
				{Label: "token", Value: token},
			},
		}, defaultToken)

	integrationResponse := structs.PipelineIntegrationResponse{}
	_ = json.Unmarshal([]byte(result), &integrationResponse)

	return integrationResponse.ID
}

func AddPipelineSource(pipelineServer string, integrationId int, repoName string, source string) {
	utils.InvokeRequestWithToken("http://"+pipelineServer+":30000/pipelineSources", "POST",
		structs.PipelineSource{
			ProjectID:            1,
			Branch:               "master",
			FileFilter:           source,
			ProjectIntegrationID: integrationId,
			RepositoryFullName:   repoName,
		}, defaultToken)
}

func InstallPipelineChart(setupInfo structs.SetupInfo, apiAddress string, wwwAddress string, artiAddress string) {
	pipelineFile := setupInfo.TempDir + "pipelines/values-ingress.yaml"
	genPipelineFile := setupInfo.TempDir + "pipelines/values-ingress.yaml"

	utils.CopyResourceToLocation("pipelines/values-ingress.yaml", pipelineFile, utils.GetAgentResources())

	pipelineChart := structs.PipelineChart{}
	utils.LoadYamlFile(pipelineFile, &pipelineChart)

	pipelineChart.API.Artifactory.URL = "http://" + artiAddress

	pipelineChart.API.Service.LoadBalancerIP = apiAddress
	pipelineChart.API.ExternalURL = "http://" + apiAddress + ":30000"

	pipelineChart.Www.Service.LoadBalancerIP = wwwAddress
	pipelineChart.Www.ExternalURL = "http://" + wwwAddress + ":30001"

	chartYaml, _ := yaml.Marshal(pipelineChart)
	_ = ioutil.WriteFile(genPipelineFile, []byte(chartYaml), 0644)

	utils.HelmInstall(setupInfo.Cluster.Name, "entplusrepo/pipelines", "", nil, genPipelineFile)
}
