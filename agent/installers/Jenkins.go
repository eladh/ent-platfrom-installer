package installers

import (
	"agent/utils"
	"bytes"
	"common/kubernetes"
	commonStructs "common/structs"
	commonUtils "common/utils"
	"encoding/base64"
	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// todo : so ugly - use padding instead
const paramSpacing string = "         "

func InstallJenkins(setupInfo commonStructs.SetupInfo) {
	logger.Always("Install Jenkins Server")

	if setupInfo.Tools.Jenkins.Site == "" {
		return
	}

	configs := make(map[interface{}]interface{})
	root := make(map[interface{}]interface{})
	env := make(map[string]string)
	image := make(map[string]string)

	image["tag"] = setupInfo.Services.Versions.Jenkins

	artifactoryConf := commonUtils.GetYamlResourceAsMap("jenkins/artifactory.yml", commonUtils.GetAgentResources())
	dslScriptsConf := commonUtils.GetYamlResourceAsMap("jenkins/dsl-scripts.yml", commonUtils.GetAgentResources())
	ldapConf := commonUtils.GetYamlResourceAsMap("jenkins/ldap.yml", commonUtils.GetAgentResources())
	securityConf := commonUtils.GetYamlResourceAsMap("jenkins/security.yml", commonUtils.GetAgentResources())
	toolsConf := commonUtils.GetYamlResourceAsMap("jenkins/tools.yml", commonUtils.GetAgentResources())
	sonarConfig := commonStructs.SonarConfig{}

	var credentialsConf map[interface{}]interface{}

	if commonUtils.IsResourceExists("jenkins/credentials.yml", commonUtils.GetAgentResources()) {
		credentialsConf = commonUtils.GetYamlResourceAsMap("jenkins/credentials.yml", commonUtils.GetAgentResources())
	} else {
		credentialsConf = commonUtils.GetYamlResourceAsMap("jenkins/default-credentials.yml", commonUtils.GetAgentResources())
	}

	certificate, _ := ioutil.ReadFile(setupInfo.TempDir + "ca.crt")
	artifactoryKey := credentialsConf["credentials"].(map[interface{}]interface{})["artipublickey"]
	artifactoryKey.(map[interface{}]interface{})["text"] = base64.StdEncoding.EncodeToString(certificate)

	configs["credentials"] = credentialsConf["credentials"]

	addSonarSupport(sonarConfig, setupInfo)

	env["ARTIFACTORY_URL"] = "http://" + GetAuthServerFullName(GetAuthProvider(setupInfo).Name, setupInfo)
	env["JENKINS_ADMIN_PASSWORD"] = "password"

	configs["artifactory"] = artifactoryConf["artifactory"]
	configs["job_dsl_scripts"] = dslScriptsConf["job_dsl_scripts"]
	configs["script_approval"] = dslScriptsConf["script_approval"]
	configs["permissions"] = ldapConf["permissions"]
	configs["tools"] = toolsConf["tools"]
	configs["sonar_qube_servers"] = sonarConfig.SonarQubeServers

	//todo - load LDAP vs Basic security
	//configs["security"] = ldapConf["security"]
	configs["security"] = securityConf["security"]

	root["env"] = env
	root["managedConfig"] = configs
	root["image"] = image

	configB, _ := yaml.Marshal(root)

	configString := strings.Replace(string(configB), "scriptTemplate", generatePipelines(setupInfo), 1)
	_ = ioutil.WriteFile(setupInfo.TempDir+"config.yml", []byte(configString), 0644)

	commonUtils.HelmRepoAdd("odavid", "https://odavid.github.io/k8s-helm-charts")
	commonUtils.HelmInstall("jenkins", "odavid/my-bloody-jenkins", "", nil, setupInfo.TempDir+"config.yml")

	addJenkinsService(setupInfo.Tools.Jenkins.Site)
}

func addSonarSupport(sonarConfig commonStructs.SonarConfig, setupInfo commonStructs.SetupInfo) {
	commonUtils.GetYamlResource("jenkins/sonar.yml", &sonarConfig, commonUtils.GetAgentResources())
	if setupInfo.Tools.Sonarqube {
		sonarConfig.SonarQubeServers.Installations.MySonarQube.ServerURL = GetSonarQubeAddress()
		sonarConfig.SonarQubeServers.Installations.MySonarQube.ServerAuthenticationToken = GenerateSonarQubeToken()
	}
}

func generatePipelines(setupInfo commonStructs.SetupInfo) string {
	// todo - replace with template engine
	var dslScript bytes.Buffer

	for _, job := range setupInfo.Tools.Jenkins.Jobs {
		script := commonUtils.GetResource("jenkins/pipeline-template.yml", commonUtils.GetAgentResources())
		script = strings.Replace(script, "job", job.Name+"|"+job.Url+"|"+job.Pipeline, 1)
		script = strings.Replace(script, "customParams", generateParams(job), 1)
		dslScript.WriteString(script + "\n")
	}

	return dslScript.String()
}

func generateParams(job commonStructs.Job) string {
	var paramsBuffer bytes.Buffer

	paramsBuffer.WriteString("parameters { \n")

	if len(job.Params) == 0 {
		return ""
	}

	for _, param := range job.Params {
		paramsBuffer.WriteString(paramSpacing + param.Type + "('" + param.Name + "'," + param.DefaultValue + ", '" + param.Desc + "')")
	}
	paramsBuffer.WriteString("} \n")

	return paramsBuffer.String()
}

func addJenkinsService(site string) {
	jfmc := kubernetes.GetServiceIp("mission-control")
	jenkins := kubernetes.GetServiceIp("jenkins-my-bloody-jenkins") + ":8080"

	utils.WaitForService("jenkins-my-bloody-jenkins", 8080, false)

	commonUtils.InvokeRequest("http://"+jfmc+"/api/v3/services", "POST", commonStructs.AddService{
		Name:                 "Jenkins",
		Description:          "Jenkins service",
		UserName:             "admin",
		Password:             "password",
		URL:                  "http://" + jenkins,
		SiteName:             site,
		PairWithAuthProvider: true,
		AuthProvider:         false,
		Type:                 "JENKINS",
	})
}
