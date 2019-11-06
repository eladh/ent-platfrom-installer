package utils

import (
	"strings"
)

func HelmRepoUpdate()  {
	Shell("helm repo update")
}

func HelmRepoAdd(alias string ,location string)  {
	Shell("helm repo add " + alias + " " + location)
}

func HelmUpgrade(name string ,repo string ,version string ,params []string)  {
	helmInternal("upgrade" ,name ,repo ,version ,params ,"","")
}

func HelmInstall(name string ,repo string ,version string ,params []string ,paramFile string)  {
	helmInternal("install" ,name ,repo ,version ,params ,paramFile ,"")
}

func HelmInstallLocalChart(name string ,path string) {
	helmInternal("install" ,name ,"" ,"" ,nil , "" ,path)
}

func helmInternal(command string ,name string ,repo string ,version string ,params []string ,paramFile string ,location string)  {
	var helmInstallCmd strings.Builder

	helmInstallCmd.WriteString("helm tiller run -- helm " + command + " ")

	if command == "install" {
		helmInstallCmd.WriteString(" --name " + name + " ")
	} else {
		helmInstallCmd.WriteString(name + " ")
	}

	if repo != "" {
		helmInstallCmd.WriteString(repo)
	}

	if version != "" {
		helmInstallCmd.WriteString(" --version " + version)
	}

	// todo - need refactor to use --set for each
	if params != nil {
		helmInstallCmd.WriteString(" " + strings.Join(params[:], " "))
	}

	if paramFile != "" {
		helmInstallCmd.WriteString(" -f " + paramFile)
	}

	if location != "" {
		helmInstallCmd.WriteString(" " + location)
	}

	Shell(helmInstallCmd.String())
}