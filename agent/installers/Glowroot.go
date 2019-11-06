package installers

import (
	agentUtils "agent/utils"
	"common/kubernetes"
	"common/structs"
	"common/utils"
	"github.com/kris-nova/logger"
	"strings"
)

const glowrootJvmAgentCommand = "-javaagent:/var/opt/jfrog/artifactory/glowroot/glowroot.jar"
const glowrootBintrayUrl = "https://github.com/glowroot/glowroot/releases/download/v0.13.4/glowroot-0.13.4-dist.zip"
const glowrootZipFile = "glowroot-0.13.4-dist.zip"
const glowrootServiceName = "artifactory-apm-glowroot"
const glowrootPropertiesFile = "/var/opt/jfrog/artifactory/glowroot/glowroot.properties"

func InstallGlowRoot(setupInfo structs.SetupInfo) {
	if !setupInfo.Tools.Glowroot {
		return
	}

	logger.Always("Install GlowRoot Server")

	utils.HelmInstallLocalChart("artifactory-apm", utils.GetResourceLocation("charts/glowroot", utils.GetAgentResources()))
	agentUtils.WaitForService(glowrootServiceName, 80, false)

	registerGlowrootAgents(setupInfo, GetArtifactoryServers(setupInfo))
}

func registerGlowrootAgents(setupInfo structs.SetupInfo, servers []structs.Artifactory) {
	file, _ := utils.DownloadFile(setupInfo.TempDir+glowrootZipFile, glowrootBintrayUrl)
	glowrootAgentAddress := kubernetes.GetServiceIp(glowrootServiceName) + ":8181"

	for _, server := range servers {
		AddGlowrootArtifactoryAgent(server.Name, file, glowrootAgentAddress)
		UpdateArtifactoryAttribute(server, setupInfo, "artifactory.javaOpts.other", glowrootJvmAgentCommand)
	}
}

func AddGlowrootArtifactoryAgent(serviceName string, glowrootContent []byte, glowrootAddress string) {
	pod, _ := kubernetes.GetPod(serviceName + "-artifactory-0")

	_, _, _ = kubernetes.UploadContentToPod(glowrootContent, "/tmp/"+glowrootZipFile, pod.Name, "artifactory", "default")
	_, _, _ = kubernetes.ExecOnPod("unzip /tmp/"+glowrootZipFile+" -d  /var/opt/jfrog/artifactory", pod.Name, "artifactory", "default", nil)

	_, _, _ = kubernetes.UploadContentToPod([]byte(generateGlowrootProperties(serviceName, glowrootAddress)),
		glowrootPropertiesFile, pod.Name, "artifactory", "default")
}

func generateGlowrootProperties(agentId string, collectorAddress string) string {
	var b1 strings.Builder
	b1.WriteString("agent.id=" + agentId + "\n")
	b1.WriteString("collector.address=" + collectorAddress + "\n")

	return b1.String()
}
