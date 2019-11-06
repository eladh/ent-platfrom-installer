package installers

import (
	"common/utils"
	"github.com/kris-nova/logger"
	"time"
)

func InstallHelmClient() {
	logger.Always("Install Helm + Tiller")

	utils.ShellNoExit("helm init --client-only")
	utils.ShellNoExit("helm plugin install https://github.com/rimusz/helm-tiller --version=0.9.0")
	
	time.Sleep(60 * 1 * time.Second)

	utils.HelmRepoAdd("jfrog" , "https://charts.jfrog.io/")
	utils.HelmRepoUpdate()
}