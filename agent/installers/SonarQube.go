package installers

import (
	agentUtils "agent/utils"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
)

func GenerateSonarQubeToken() string {
	server := kubernetes.GetServiceIp("sonarqube-server-sonarqube")
	tokenJson := commonUtils.InvokeRequestWithPassword("http://"+server+":9000/api/user_tokens/generate", "admin", "admin",
		"POST", "name=admin")

	return commonUtils.GetJsonStringAttribute(".token", tokenJson)
}

func InstallSonarQube(setupInfo structs.SetupInfo) {
	if !setupInfo.Tools.Sonarqube {
		return
	}

	logger.Always("Install SonarQube Server")
	commonUtils.
		HelmInstall("sonarqube-server" ,"stable/sonarqube" ,setupInfo.Services.Versions.Sonar ,[]string {} ,"")
	agentUtils.WaitForService("sonarqube-server-sonarqube", 9000, true)
}

func GetSonarQubeAddress() string {
	return "http://" + kubernetes.GetServiceIp("sonarqube-server-sonarqube") + ":9000"
}
