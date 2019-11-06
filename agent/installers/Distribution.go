package installers

import (
	agentUtils "agent/utils"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
)

func GeneratePgpKeys(servers []structs.Artifactory ,setupInfo structs.SetupInfo) {
	logger.Always("Setup Distribution + Artifactory  gpg keys")

	if setupInfo.Services.Distribution.Site == "" {
		return
	}

	ip := kubernetes.GetServiceIp(setupInfo.Services.Distribution.Name)

	commonUtils.CopyResourceToLocation("gpg/gen-key-params" , setupInfo.TempDir + "/gen-key-params" ,commonUtils.GetAgentResources())

	commonUtils.Shell("gpg --pinentry-mode=loopback --passphrase \"\" --quiet --batch --no-tty" +
		" --gen-key	" + setupInfo.TempDir + "/gen-key-params")

	//Find the latest key
	id := commonUtils.Shell("gpg --no-tty --list-secret-keys --with-colons 2>/dev/null | " +
		"tail | awk -F: '/^sec:/ { print $5 } ' | tail -n 1")

	privateKey := commonUtils.Shell("gpg --armor --export-secret-keys " + id)
	publicKey := commonUtils.Shell("gpg --armor --export " + id)

	commonUtils.InvokeRequest("http://"+ip+"/api/v1/keys/pgp", "PUT", structs.UpdateDistributionKeysPayload{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	})

	for _, artifactory := range servers {
		serviceIp := kubernetes.GetServiceIp(artifactory.Name + artifactoryNginxService)
		commonUtils.InvokeRequest("http://"+serviceIp+"/artifactory/api/security/keys/trusted", "POST", structs.AddArtifactoryTrustedKeyPayout{
			PublicKey: publicKey,
			Alias: artifactory.Name +"Key",
		})
	}
}

func InstallDistributionServer(setupInfo structs.SetupInfo) {
	if setupInfo.Services.Distribution.Name == "" {
		return
	}

	logger.Always("Install Distribution Server")

	distributionServerName := setupInfo.Services.Distribution.Name


	commonUtils.HelmInstall(distributionServerName , "jfrog/distribution" , setupInfo.Services.Versions.Distribution,
		[]string {setDistributionLoggers(),},"")

	agentUtils.WaitForService(distributionServerName, 80, true)
}

func AddDistributionServices(setupInfo structs.SetupInfo) {
	logger.Always("Setup Distribution service")

	if setupInfo.Services.Distribution.Site == "" {
		return
	}

	siteName := setupInfo.Services.Distribution.Site
	ip := kubernetes.GetServiceIp("mission-control")
	distributionIp := kubernetes.GetServiceIp("distribution")

	commonUtils.InvokeRequest("http://"+ip+"/api/v3/services", "POST", structs.AddService{
		Name:                 "distribution",
		Description:          "main distribution service",
		URL:                  "http://" + distributionIp,
		SiteName:             siteName,
		PairWithAuthProvider: true,
		AuthProvider:         false,
		Type:                 "DISTRIBUTION",
	})
}

func setDistributionLoggers() string {
	return " --set distributor.loggers[0]=distributor.log " +
		" --set distributor.loggers[1]=foreman.log " +
		" --set distribution.loggers[0]=access.log " +
		" --set distribution.loggers[1]=distribution.log " +
		" --set distribution.loggers[2]=request.log "
}
