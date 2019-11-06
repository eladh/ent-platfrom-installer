package main

import (
	"common/cloud/gcp"
	"common/structs"
	"common/utils"
	"github.com/kris-nova/logger"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

var setupInfoLocation string
var secretLocation string
var installerVersion string
var clusterPrefix string
var entBucketLocation string
var edgeBucketLocation string
var jenkinsSecretsLocation string
var instances int

const bintrayAgentArtifactUrl string = "jfrog-int-docker-eplus-installer-images.bintray.io/agent/core"

func main() {
	app := cli.NewApp()
	app.Name = "Platform Installer"
	app.Usage = "make infra as code E+ automation"
	app.Version = "1.0.0"

	var installCommandFlags = []cli.Flag {
		cli.StringFlag{
			Name: "config",
			Usage: "Load setupInfo configuration from file",
			Destination: &setupInfoLocation,
		},
		cli.StringFlag{
			Name: "secret",
			Usage: "Load secret configuration from file",
			Destination: &secretLocation,
		},
		cli.StringFlag{
			Name: "version",
			Usage: "Load version configuration",
			Destination: &installerVersion,
		},
		cli.StringFlag{
			Name: "name",
			Usage: "set cluster names prefix",
			Destination: &clusterPrefix,
		},
		cli.IntFlag{
			Name: "instances",
			Usage: "Load num instances`",
			Destination: &instances,
		},
		cli.StringFlag{
			Name: "artifactory",
			Usage: "Load ent bucket location",
			Destination: &entBucketLocation,
		},
		cli.StringFlag{
			Name: "edge",
			Usage: "set edge bucket location",
			Destination: &edgeBucketLocation,
		},
		cli.StringFlag{
			Name: "jenkins-secrets",
			Usage: "set jenkins secret location",
			Destination: &jenkinsSecretsLocation,
		},
	}

	var listCommandFlags = []cli.Flag {
		cli.StringFlag{
			Name: "secret",
			Usage: "Load secret configuration from file",
			Destination: &secretLocation,
		},
		cli.StringFlag{
			Name: "name",
			Usage: "set cluster names prefix",
			Destination: &clusterPrefix,
		},
		cli.IntFlag{
			Name: "instances",
			Usage: "Load num instances`",
			Destination: &instances,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install platform",
			Flags: installCommandFlags,
			Action: func(c *cli.Context) error {
				err := validate()
				if err != nil {
					return err
				}
				err = validateInstall()
				if err != nil {
					return err
				}
				ParallelDockerCommandRunner("install" ,setupInfoLocation ,secretLocation ,installerVersion ,clusterPrefix , instances ,entBucketLocation ,edgeBucketLocation ,jenkinsSecretsLocation)
				return nil
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "Uninstall platform",
			Flags: installCommandFlags,
			Action: func(c *cli.Context) error {
				err := validate()
				if err != nil {
					return err
				}
				ParallelDockerCommandRunner("uninstall" ,setupInfoLocation ,secretLocation ,installerVersion ,clusterPrefix , instances ,"" ,"" ,"")
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List platform SSH servers",
			Flags: listCommandFlags,
			Action: func(c *cli.Context) error {
				if secretLocation == "" {
					return cli.NewExitError("you must supply setupInfo yaml file by using --secret file`", 86)
				}
				if clusterPrefix == "" {
					return cli.NewExitError("you must supply clusterPrefix by using --name value", 86)
				}
				if instances == 0 {
					return cli.NewExitError("you must supply instances by using --instances value", 86)
				}

				GetSshServerAddress(secretLocation ,clusterPrefix , instances)
				return nil
			},
		},

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func validate() *cli.ExitError {
	if setupInfoLocation == "" {
		return cli.NewExitError("you must supply setupInfo yaml file by using --config file`", 86)
	}
	if secretLocation == "" {
		return cli.NewExitError("you must supply secret file by using --secret file`", 86)
	}
	if installerVersion == "" {
		return cli.NewExitError("you must supply version by using --version value", 86)
	}
	if clusterPrefix == "" {
		return cli.NewExitError("you must supply clusterPrefix by using --name value", 86)
	}
	if instances == 0 {
		return cli.NewExitError("you must supply instances by using --instances value", 86)
	}

	return nil
}

func validateInstall() *cli.ExitError {
	if entBucketLocation == "" {
		return cli.NewExitError("you must supply ent Bucket Location by using --artifactory bucket location`", 86)
	}
	if edgeBucketLocation == "" {
		return cli.NewExitError("you must supply edge Bucket Location by using --edge bucket location	`", 86)
	}

	return nil
}

func ParallelDockerCommandRunner(cmd string,setupInfoLocation string ,licenseLocation string ,version string ,clusterPrefix string ,instances int ,entBucketLocation string ,edgeBucketLocation string ,jenkinsSecretsLocation string) {
	logger.Always("start platform install")
	var myWaitGroup sync.WaitGroup

	myWaitGroup.Add(instances)

	for i := 1; i <=instances; i++ {
		go runDockerCommand(cmd ,&myWaitGroup , bintrayAgentArtifactUrl + ":" + version ,clusterPrefix, setupInfoLocation, i, licenseLocation ,entBucketLocation ,edgeBucketLocation ,jenkinsSecretsLocation)
	}

	myWaitGroup.Wait()
}


func loadSetupInfo(setupFileLocation string) structs.SetupInfo {
	setupInfo := structs.SetupInfo{}
	utils.LoadYamlFile(setupFileLocation , &setupInfo)
	setupInfo.TempDir = "/home/appuser/temp/"

	return setupInfo
}

func runDockerCommand(cmd string ,myWaitGroup *sync.WaitGroup ,dockerImageTag string ,clusterPrefix string ,setupInfoLocation string, instanceId int, licenseLocation string ,entBucketLocation string ,edgeBucketLocation string ,jenkinsSecretsLocation string) {
	setupInfo := loadSetupInfo(setupInfoLocation)
	clusterName := clusterPrefix + "-" + strconv.Itoa(instanceId)
	dir ,err:=os.Getwd()
	dir+="/temp/"

	if err != nil {
		logger.Critical("error")
	}

	_ = os.MkdirAll(dir, 0777)

	setupInfoTemp := dir + clusterName + ".yaml"
	setupInfoTempLog := dir + clusterName + ".log"
	setupInfo.Cluster.Name = clusterName
	setupInfo.Cluster.Domain = clusterName
	setupInfoYaml, _ := yaml.Marshal(setupInfo)

	ioutil.WriteFile(setupInfoTemp, setupInfoYaml, 0777)

	dockerCmd := "docker run " +
		" -v " + licenseLocation  + ":/home/appuser/license.json " +
		" -v " + setupInfoTemp    + ":/setup.yaml "

	if entBucketLocation != "" {
		dockerCmd += " -v " + entBucketLocation  + ":/home/appuser/agent/resources/buckets/artifactory.json ";
	}

	if edgeBucketLocation != "" {
		dockerCmd += " -v " + edgeBucketLocation + ":/home/appuser/agent/resources/buckets/edge.json ";
	}

	if jenkinsSecretsLocation != "" {
		dockerCmd += " -v " + jenkinsSecretsLocation + ":/home/appuser/agent/resources/jenkins/credentials.yml ";
	}

	utils.Shell(dockerCmd + " -e COMMAND=" + cmd + " " + dockerImageTag + " > " + setupInfoTempLog )

	myWaitGroup.Done()
}

func GetSshServerAddress(licenseLocation string ,name string ,instances int) {
	gcp.LoginByServiceAccount(licenseLocation)

	var data [][]string
	for i := 1; i < instances + 1; i++ {
		clusterName :=  name + "-" + strconv.Itoa(i) + "-cluster";
		gcp.Connect2Cluster(clusterName ,"devops-consulting")
		ip := utils.Shell("kubectl get svc | grep  dev-sshd-dev | awk '{print $4}' | head -n 1")
		data = append(data, []string{clusterName, ip})
	}
	printTable(data)
}

func printTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cluster name", "IP"})
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()
}