package installers

import (
	"agent/utils"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
	_ "io"
	"io/ioutil"
	"log"
	"net/http"
)

const xrayThreatsUrl = "https://www.googleapis.com/storage/v1/b/xray-storage/o/vuln_1550707199999.zip?alt=media"

func getXrayToken(address string) string {
	request, _ := commonUtils.LoadJsonAsset(structs.LoginPayload{
		Name:     "admin",
		Password: "password",
	}, "POST", "http://"+address+"/ui/auth/login" ,"application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	token, _ := ioutil.ReadAll(resp.Body)

	return "Bearer " + string(commonUtils.GetJsonStringAttribute(".token", string(token)))
}

func InstallXrayServers(setupInfo structs.SetupInfo) {
	logger.Always("Install Xray Servers")

	for _, server := range setupInfo.Services.Xray {
		commonUtils.HelmInstall(server.Name ,"jfrog/xray" ,setupInfo.Services.Versions.Xray,
			[]string {"--set common.xrayConfig.indexAllBuilds=true"} ,"")

		utils.WaitForService(server.Name, 80, true)
	}
}

func SetupXrayServers(setupInfo structs.SetupInfo) {
	logger.Always("Setup Xray servers")

	for _, server := range setupInfo.Services.Xray {
		xrayAddress := kubernetes.GetServiceIp(server.Name)
		artifactoryAddress := kubernetes.GetServiceIp(server.Artifactory + artifactoryNginxService)

		finishInstallWizard(xrayAddress)
		setArtifactoryIndexer(xrayAddress, artifactoryAddress, server.Artifactory)
		indexedAssets(xrayAddress, server, "artifactory")
		AddWatchesAndPolicies(xrayAddress, server)
		AddXrayServices(server, xrayAddress)
	}
}

func SyncXrayVulnerabilities(setupInfo structs.SetupInfo) {
	logger.Always("Download  Vulnerabilities file")
	file  ,_ := commonUtils.DownloadFile(setupInfo.TempDir+"vuln.zip", xrayThreatsUrl)

	for _, server := range setupInfo.Services.Xray {
		DownloadVulnerabilities(server.Name, file)
	}
}

func DownloadVulnerabilities(xrayServer string, file []byte) {
	pod, _ := kubernetes.GetPod(xrayServer)
	logger.Always("Upload  Vulnerabilities file to pod " + xrayServer)
	_, _, _ = kubernetes.UploadContentToPod(file, "/var/opt/jfrog/xray/data/updates/vulnerability/vuln.zip", pod.Name, "", "default")
}

func AddWatchesAndPolicies(ip string, server structs.XrayServer) {
	commonUtils.InvokeRequest("http://"+ip+"/api/v1/policies", "POST", server.Policies[0])
	commonUtils.InvokeRequest("http://"+ip+"/api/v2/watches", "POST", server.Watches[0])
}

func indexedAssets(ip string, server structs.XrayServer, artifactoryName string) {
	commonUtils.InvokeRequest("http://"+ip+"/api/v1/binMgr/"+artifactoryName+"/repos", "PUT",
		structs.UpdateIndexedRepos{Repos: server.Repos})

	commonUtils.InvokeRequest("http://"+ip+"/api/v1/binMgr/"+artifactoryName+"/builds", "PUT",
		structs.UpdateIndexedBuilds{Builds: server.Builds})
}

func setArtifactoryIndexer(xrayAddress string, artifactoryAddress string, artifactoryServiceName string) {
	commonUtils.InvokeRequestWithToken("http://"+xrayAddress+"/ui/binMgr", "POST", structs.IndexedArtifactoryPayload{
		User:         "admin",
		Password:     "password",
		BinMgrURL:    "http://" + artifactoryAddress + "/artifactory",
		BinMgrID:     artifactoryServiceName,
		BinMgrDesc:   artifactoryServiceName,
		ProxyEnabled: false,
	}, getXrayToken(xrayAddress))
}

func finishInstallWizard(ip string) {
	commonUtils.InvokeRequestWithToken("http://"+ip+"/ui/onboardingconfig", "POST", structs.WizardPayload{
		BaseURL:  "http://" + ip,
		Finished: true,
		StepNum:  5,
	}, getXrayToken(ip))
}

func AddXrayServices(server structs.XrayServer, address string) {
	ip := kubernetes.GetServiceIp("mission-control")

	commonUtils.InvokeRequest("http://"+ip+"/api/v3/services", "POST", structs.AddService{
		Name:                 server.Name,
		Description:          "xray " + server.Name + " service",
		UserName:             "admin",
		Password:             "password",
		URL:                  "http://" + address,
		SiteName:             server.Site,
		PairWithAuthProvider: true,
		AuthProvider:         false,
		Type:                 "XRAY",
	})
}
