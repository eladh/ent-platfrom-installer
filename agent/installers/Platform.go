package installers

import (
	agentUtils "agent/utils"
	"common/cloud/aws"
	"common/cloud/gcp"
	"common/installers"
	"common/kubernetes"
	"common/structs"
	"common/utils"
	"github.com/kris-nova/logger"
	"log"
	"os"
)


func InstallPlatform(setupInfoLocation string) {
	logger.Always("start platform install")

	initLogger()
	setupInfo := loadSetupInfo(setupInfoLocation)

	generateTempDir(setupInfo.TempDir)

	installKubernetesCluster(setupInfo)
	installers.InstallHelmClient()
	installServices(setupInfo)
	setupServices(setupInfo)
}

func UninstallPlatform(setupInfoLocation string) {
	logger.Always("start platform uninstall")

	setupInfo := loadSetupInfo(setupInfoLocation)
	if setupInfo.Vendor.Type == "gcp" {
		gcp.Uninstall(setupInfo)
	} else {
		log.Panic("unable to install platform on target env")
		//	aws.Uninstall(setupInfo)
	}
}

func InstallArtifactory(name string, site string ,setupInfoLocation string) {
	AddArtifactoryInstance(loadSetupInfo(setupInfoLocation), structs.Artifactory{Name: name, Site: site})
	AddRepositories([]structs.Artifactory{{Site: site, Name: name, AuthServer: false}})
}

func installKubernetesCluster(setupInfo structs.SetupInfo) {
	if setupInfo.Vendor.Type == "gcp" {
		gcp.Install(setupInfo)
	} else  if setupInfo.Vendor.Type == "aws" {
		aws.Install(setupInfo)
	} else {
		log.Panic("unable to install platform on target env")
	}
}

func installServices(setupInfo structs.SetupInfo) {
	logger.Always("Installing Services")

	InstallMissionControl("mission-control", setupInfo)
	InstallArtifactoryServers(setupInfo)
	InstallArtifactoryEdgeServers(setupInfo)
	InstallDistributionServer(setupInfo)
	InstallXrayServers(setupInfo)
	InstallSonarQube(setupInfo)

	if setupInfo.Tools.Dev {
		utils.HelmInstallLocalChart("dev",utils.GetResourceLocation("charts/dev/chart" ,utils.GetAgentResources()))
	}
	agentUtils.WaitForAllServices()
}


func setupServices(setupInfo structs.SetupInfo) {
	logger.Always("Setup Services")

	artifactoryServers := GetArtifactoryServers(setupInfo)

	if !setupInfo.Services.InstallOnly {
		SetupMissionControl(setupInfo)
		SetupArtifactoryServers(setupInfo)
		SetupEdgeServers(setupInfo)
		AddRepositories(artifactoryServers)
		SetupCircleOfTrust(artifactoryServers)
		AddFederationRules(artifactoryServers)
		SetupReplications(artifactoryServers)
		AddDistributionServices(setupInfo)
		GeneratePgpKeys(artifactoryServers, setupInfo)
		SetupXrayServers(setupInfo)
	}

	if setupInfo.Tools.Dev {
		login2Cluster(setupInfo)
	}


	if setupInfo.Tools.Glowroot {
		InstallGlowRoot(setupInfo)
	}

	SyncXrayVulnerabilities(setupInfo)
	InstallJenkins(setupInfo)
}

func login2Cluster(setupInfo structs.SetupInfo) {
	sshPod, _ := kubernetes.GetPod("dev-sshd-dev")
	loginScript := utils.GetResource("ssh/login.sh", utils.GetAgentResources())
	_, _, _ = kubernetes.UploadFileToPod("/home/appuser/license.json", "/tmp/license.json", sshPod.Name, "sshd-dev", "default")
	_, _, _ = kubernetes.UploadContentToPod([]byte(loginScript), "/tmp/login.sh", sshPod.Name, "sshd-dev", "default")
	_, _, _ = kubernetes.ExecOnPod("chmod +x /tmp/login.sh", sshPod.Name, "sshd-dev", "default", nil)
	_, _, _ = kubernetes.ExecOnPod("/tmp/login.sh "+setupInfo.Cluster.Name+"-cluster", sshPod.Name, "sshd-dev", "default", nil)
}

func loadSetupInfo(setupFileLocation string) structs.SetupInfo {
	logger.Always("Use configuration file = "+ setupFileLocation)

	if setupFileLocation == "" || !utils.FileExists(setupFileLocation) {
		logger.Critical("Please setup proper setupInfo yaml file location")
		os.Exit(1)
	}

	setupInfo := structs.SetupInfo{}
	utils.LoadYamlFile(setupFileLocation , &setupInfo)

	return setupInfo
}

func initLogger() {
	logger.Fabulous = false
	logger.Color = false
	logger.Level = 4
}

func generateTempDir(tempDir string) {
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0777)
	logger.Always("temp dir is " + tempDir)
}
