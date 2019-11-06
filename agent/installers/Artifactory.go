package installers

import (
	agentUtils "agent/utils"
	"bufio"
	"common/cloud/gcp"
	"common/eventbus"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const artifactoryNginxService = "-artifactory-nginx"

func AddArtifactoryInstance(setupInfo structs.SetupInfo, artifactory structs.Artifactory) {
	addArtifactory(setupInfo, artifactory)
	agentUtils.WaitForAllServices()

	artifactoryAddress := kubernetes.GetServiceIp(artifactory.Name + artifactoryNginxService)
	addArtifactoryLicense(artifactory.Name, artifactoryAddress, 80, "enterprise-buckets")
	addArtifactoryService(artifactory.Name, artifactoryAddress, 80, artifactory.Site, false)
}

func InstallArtifactoryServers(setupInfo structs.SetupInfo) {
	logger.Always("Install Artifactory Servers")

	for _, server := range setupInfo.Services.Artifactory {
		installArtifactory(server, setupInfo)
	}
}

func InstallArtifactoryEdgeServers(setupInfo structs.SetupInfo) {
	logger.Always("Install Artifactory Edge Servers")

	for _, server := range setupInfo.Services.Edges {
		_ = kubernetes.CreateNamespace(server.Name)
		installEdgeArtifactory(server.Name, setupInfo.Services.Versions.Artifactory)
	}
}

func installArtifactory(artifactory structs.Artifactory, setupInfo structs.SetupInfo) {
	if artifactory.HA != (structs.HighAvailability{}) {
		logger.Always("Install Artifactory HA settings")
		addArtifactoryHA(setupInfo, artifactory)
	}

	if artifactory.AuthServer {
		addArtifactoryWithSslSupport(artifactory, GetAuthServerFullName(artifactory.Name, setupInfo), setupInfo)
	} else {
		addArtifactory(setupInfo, artifactory)
	}
}

func GetAuthServerFullName(serverName string, setupInfo structs.SetupInfo) string {
	domain := setupInfo.Cluster.Domain
	cluster := setupInfo.Cluster.Name

	return serverName + "." + cluster + "." + domain
}

func SetupArtifactoryServers(setupInfo structs.SetupInfo) {
	logger.Always("Setup Artifactory servers")

	for _, artifactory := range setupInfo.Services.Artifactory {
		artifactoryAddress := kubernetes.GetServiceIp(artifactory.Name + artifactoryNginxService)

		if artifactory.AuthServer {
			InstallPipeline(setupInfo, artifactoryAddress)
			//todo - need to be cloud agnostic
			gcp.CreateInternalDnsRecords(setupInfo, artifactory.Name)
		}

		if artifactory.HA != (structs.HighAvailability{}) {
			logger.Always("Setup Artifactory HA settings")

			waitForHACluster(artifactory.Name+"-nginx", 80)
			addArtifactoryLicense(artifactory.Name, artifactoryAddress, 80, "enterprise-licenses")
			addArtifactoryService(artifactory.Name, artifactoryAddress, 80, artifactory.Site, artifactory.AuthServer)
		} else {
			addArtifactoryLicense(artifactory.Name, artifactoryAddress, 80, "enterprise-licenses")
			addArtifactoryService(artifactory.Name, artifactoryAddress, 80, artifactory.Site, artifactory.AuthServer)
		}
	}
}

func UpdateArtifactoryVersion(server structs.Artifactory, setupInfo structs.SetupInfo) {
	latestVersion := agentUtils.GetProductsVersions()["artifactory"][0]
	UpdateArtifactoryAttribute(server, setupInfo, "artifactory.image.version", latestVersion.Original())
}

func UpdateArtifactoryAttribute(server structs.Artifactory, setupInfo structs.SetupInfo, key string, value string) {
	params := []string{" --set " + key + "=" + value, " --set postgresql.postgresPassword=zooloo"}
	if server.AuthServer {
		tlsParams := generateTlsParams(server.Name)
		params = append(params, tlsParams[0])
		params = append(params, tlsParams[1])
	}

	commonUtils.HelmUpgrade(server.Name, "jfrog/artifactory", setupInfo.Services.Versions.Artifactory, params)
	agentUtils.WaitForService(server.Name+artifactoryNginxService, 80, false)
}

func GetArtifactoryParam(server structs.Artifactory, attribute string) string {
	artifactoryAddress := kubernetes.GetServiceIp(server.Name + artifactoryNginxService)
	response := commonUtils.InvokeRequest("http://"+artifactoryAddress+"/artifactory/api/system", "GET", "")
	response = strings.Replace(response, " ", "", -1)

	var params = make(map[string]string)
	reader := bufio.NewReader(strings.NewReader(response))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if len(line) <= 0 {
			continue
		}
		value := strings.Split(string(line), "|")

		if len(value) != 2 {
			continue
		}
		params[value[0]] = value[1]
	}

	return params[attribute]
}

func waitForHACluster(name string, port int) {
	artifactory := kubernetes.GetServiceIp(name)
	responseOk := false
	for !responseOk {
		resp, err := http.Get("http://" + artifactory + ":" + strconv.Itoa(port) + "/artifactory/ui/auth/current")
		if err == nil && resp.StatusCode == 200 {
			responseOk = true
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func SetupEdgeServers(setupInfo structs.SetupInfo) {
	logger.Always("Setup Edge nodes")

	for _, artifactory := range setupInfo.Services.Edges {
		artifactoryAddress := kubernetes.GetServiceIp(artifactory.Name + artifactoryNginxService)
		addArtifactoryLicense(artifactory.Name, artifactoryAddress, 80, "edge-licenses")
		addArtifactoryService(artifactory.Name, artifactoryAddress, 80, artifactory.Site, false)
	}
}

func installEdgeArtifactory(name string, version string) {
	logger.Always("Install Edge Servers")

	commonUtils.HelmInstall(name, "jfrog/artifactory", version,
		[]string{
			" --set postgresql.postgresPassword=zooloo ",
			" --set artifactory.replicator.enabled=true ",
			" --set artifactory.replicator.publicUrl=http://localhost:6061 ",
			setArtifactoryLoggers(),
		}, "")

	agentUtils.WaitForService(name+artifactoryNginxService, 80, true)
}

func addArtifactoryLicense(name string, address string, port int, licenseType string) {
	license := getLicense(licenseType, name)
	commonUtils.InvokeRequest("http://"+address+":"+strconv.Itoa(port)+"/artifactory/api/system/licenses", "POST", structs.ArtifactoryLicense{
		LicenseKey: license,
	})
}

func getLicense(licenseType string, name string) string {
	jfmc := kubernetes.GetServiceIp("mission-control")

	response := commonUtils.InvokeRequest("http://"+jfmc+"/api/v3/attach_lic/buckets/"+licenseType, "POST",
		structs.AttachLicencePayload{
			ServiceName: name,
			Deploy:      false,
			Instances:   1,
		})

	return commonUtils.GetJsonStringUnquoteAttribute(".license_key", response)
}

func addArtifactoryService(name string, address string, port int, site string, authProvider bool) {
	jfmc := kubernetes.GetServiceIp("mission-control")

	commonUtils.InvokeRequest("http://"+jfmc+"/api/v3/services", "POST", structs.AddService{
		Name:                 name,
		Description:          name + " service",
		URL:                  "http://" + address + ":" + strconv.Itoa(port) + "/artifactory",
		SiteName:             site,
		UserName:             "admin",
		Password:             "password",
		PairWithAuthProvider: !authProvider,
		AuthProvider:         authProvider,
		Type:                 "ARTIFACTORY",
	})
}

func setTlsSupport(artifactory structs.Artifactory, fqdn string, setupInfo structs.SetupInfo) string {
	cert, key, _ := commonUtils.GenerateCert("someOrg", "someUnit", "US", "san fran", fqdn,
		[]string{fqdn, "docker." + fqdn}, setupInfo.TempDir)

	//todo - must be cloud agnostic
	gcp.DeployCertToClusterNodes(setupInfo.Vendor.Zone, "docker."+fqdn, setupInfo.TempDir)

	artifactoryName := artifactory.Name

	tlsSecret := artifactoryName + "-tls-secret"
	nginxConfig := artifactoryName + "-nginx-conf"

	_ = kubernetes.DeleteConfigMap(nginxConfig, "default")
	_ = kubernetes.DeleteSecret(tlsSecret, "default")
	_ = kubernetes.CreateSecret(tlsSecret, "default", cert, key)
	_ = kubernetes.CreateConfigMap(nginxConfig, "artifactory.conf", "default", getArtifactoryNginx(fqdn, artifactoryName))

	return strings.Join(generateTlsParams(artifactoryName)[:], " ")
}

func generateTlsParams(artifactoryName string) []string {
	tlsSecret := artifactoryName + "-tls-secret"
	nginxConfig := artifactoryName + "-nginx-conf"

	return []string{" --set nginx.tlsSecretName=" + tlsSecret, " --set nginx.customArtifactoryConfigMap=" + nginxConfig}

}

func addArtifactoryWithSslSupport(artifactory structs.Artifactory, fqdn string, setupInfo structs.SetupInfo) {
	artifactoryName := artifactory.Name

	commonUtils.HelmInstall(artifactoryName, "jfrog/artifactory", setupInfo.Services.Versions.Artifactory,
		[]string{
			" --set postgresql.persistence.size=10Gi ",
			" --set postgresql.postgresPassword=zooloo",
			setTlsSupport(artifactory, fqdn, setupInfo),
			setArtifactoryLoggers(),
			setStorageSupport(setupInfo, artifactory)}, "")

	agentUtils.WaitForService(artifactoryName+artifactoryNginxService, 80, true)
}

func addArtifactory(setupInfo structs.SetupInfo, artifactory structs.Artifactory) {
	commonUtils.HelmInstall(artifactory.Name, "jfrog/artifactory", setupInfo.Services.Versions.Artifactory,
		[]string{
			" --set postgresql.postgresPassword=zooloo ",
			" --set postgresql.persistence.size=10Gi ",
			setArtifactoryLoggers(),
			setStorageSupport(setupInfo, artifactory),
		}, "")

	agentUtils.WaitForService(artifactory.Name+artifactoryNginxService, 80, true)
}

func addArtifactoryHA(setupInfo structs.SetupInfo, artifactory structs.Artifactory) {
	commonUtils.HelmInstall(artifactory.Name, "jfrog/artifactory-ha", setupInfo.Services.Versions.ArtifactoryHA,
		[]string{
			" --set artifactory.node.replicaCount=" + strconv.Itoa(artifactory.HA.MinAvailable),
			" --set postgresql.postgresPassword=zooloo ",
			" --set postgresql.persistence.size=10Gi ",
			" --set artifactory.masterKey=" + commonUtils.GenerateRandomKey(),
			setArtifactoryLoggers(),
		}, "")

	agentUtils.WaitForService(artifactory.Name+"-nginx", 80, true)
}

func getArtifactoryNginx(fullDomainName string, artifactoryName string) string {
	artifactoryNginx := commonUtils.GetResource("artifactory/nginx.conf", commonUtils.GetAgentResources())
	artifactoryNginx = strings.Replace(artifactoryNginx, "$domain_name", fullDomainName, -1)
	artifactoryNginx = strings.Replace(artifactoryNginx, "$server_name", artifactoryName, -1)

	return artifactoryNginx
}

func GetArtifactoryServers(setupInfo structs.SetupInfo) []structs.Artifactory {
	return append(setupInfo.Services.Artifactory, setupInfo.Services.Edges...)
}

func GetAuthProvider(setupInfo structs.SetupInfo) structs.Artifactory {
	for _, server := range setupInfo.Services.Artifactory {
		if server.AuthServer {
			return server
		}
	}

	return structs.Artifactory{}
}

func setArtifactoryLoggers() string {
	return " --set artifactory.loggers[0]=request.log " +
		"--set artifactory.loggers[1]=event.log " +
		"--set artifactory.loggers[2]=access.log " +
		"--set artifactory.loggers[3]=artifactory.log " +
		"--set artifactory.catalinaLoggers[0]=catalina.log " +
		"--set nginx.catalinaLoggers[0]=access.log " +
		"--set nginx.catalinaLoggers[1]=error.log "
}

func setStorageSupport(setupInfo structs.SetupInfo, artifactory structs.Artifactory) string {
	if artifactory.Storage == (structs.Storage{}) {
		return ""
	}

	if setupInfo.Vendor.Type == "gcp" {
		bucketName := artifactory.Storage.Name + strconv.Itoa(rand.Int())
		logger.Always("create new cloud storage named " + bucketName)
		gcp.CreateGcpBucket(setupInfo.Vendor.Project, bucketName, artifactory.Storage.Type, artifactory.Storage.Location)

		return " --set artifactory.persistence.type=google-storage " +
			" --set bucketName=" + bucketName +
			" --set bucketExists=true " +
			" --set artifactory.persistence.googleStorage.identity=" + setupInfo.Vendor.Gcp.Storage.Identity +
			" --set artifactory.persistence.googleStorage.credential=" + setupInfo.Vendor.Gcp.Storage.Secret
	}

	return ""
}

func SyncArtifactoryPublicKey() {
	ch1 := make(chan eventbus.DataEvent)
	eventbus.Subscribe("k8s", ch1)

	select {
	case listen := <-ch1:
		logger.Always("Channel: %s; Topic: %s; DataEvent: %v\n", ch1, listen.Topic, listen.Data)

	}
}

func exposeHAInstances(setupInfo structs.SetupInfo) {
	for _, server := range setupInfo.Services.Artifactory {
		if server.HA != (structs.HighAvailability{}) {
			kubernetes.WaitForPod(server.Name+"-artifactory-ha-member-"+strconv.Itoa(server.HA.ReplicaCount-1), 6)
			kubernetes.ExposePod(server.Name+"-artifactory-ha-primary-0", "LoadBalancer", server.Name)

			for i := 0; i < server.HA.ReplicaCount; i++ {
				index := strconv.Itoa(i)
				kubernetes.ExposePod(server.Name+"-artifactory-ha-member-"+index, "LoadBalancer", server.Name)
			}
		}
	}
}
