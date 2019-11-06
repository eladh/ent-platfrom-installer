package gcp

import (
	"bufio"
	gcp "common/cloud/gcp/structs"
	"common/structs"
	commonUtil "common/utils"
	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func Install(setupInfo structs.SetupInfo) {
	zone := setupInfo.Vendor.Zone
	clusterName := setupInfo.Cluster.Name
	projectName := setupInfo.Vendor.Project

	gcpInit(zone, projectName)
	installCluster(clusterName, zone, setupInfo.TempDir)
	setupK8s()
}

func DeleteAllDnsRecords() {
	results := commonUtil.Shell("gcloud dns managed-zones list | awk '{print $1}'")
	scanner := bufio.NewScanner(strings.NewReader(results))
	for scanner.Scan() {
		clusterName := scanner.Text()
		if clusterName == "NAME" {
			continue
		}

		commonUtil.Shell("gcloud dns record-sets import -z " + clusterName + "  --delete-all-existing /dev/null")
		commonUtil.Shell("gcloud dns managed-zones delete " + clusterName)
	}
}

func Uninstall(setupInfo structs.SetupInfo) {
	//todo - cleanup buckets !!
	zone := setupInfo.Vendor.Zone
	clusterName := setupInfo.Cluster.Name
	projectName := setupInfo.Vendor.Project

	logger.Always("Uninstall Kubernetes Cluster")

	commonUtil.Shell("gcloud config set compute/zone " + zone)
	commonUtil.Shell("gcloud config set project " + projectName)
	commonUtil.Shell("gcloud deployment-manager deployments delete " + clusterName + "-deploy  --quiet")
	commonUtil.Shell("gcloud compute disks delete $(gcloud compute disks list | grep gke-" + clusterName + " | awk '{print $1}')")
	commonUtil.Shell("gcloud dns record-sets import -z " + clusterName + "  --delete-all-existing /dev/null")
	commonUtil.Shell("gcloud dns managed-zones delete " + clusterName)
}

func installCluster(clusterName string, zone string, tempDir string) {
	logger.Always("Install Kubernetes Cluster")

	generateInstallConfig(tempDir, clusterName ,zone)
	commonUtil.Shell("gcloud deployment-manager deployments create " + clusterName + "-deploy --config " + tempDir + "properties.yaml")
	commonUtil.Shell("gcloud container clusters get-credentials " + clusterName + "-cluster  --zone " + zone)
}

func generateInstallConfig(tempDir string, clusterName string ,zone string) {
	commonUtil.CopyResourceToLocation("cloud/gcp/deployment/properties.yaml", tempDir+"properties.yaml", commonUtil.GetCommonResources())
	commonUtil.CopyResourceToLocation("cloud/gcp/deployment/templates/cluster.jinja", tempDir+"templates/cluster.jinja", commonUtil.GetCommonResources())
	commonUtil.CopyResourceToLocation("cloud/gcp/deployment/templates/cluster.jinja.schema", tempDir+"templates/cluster.jinja.schema", commonUtil.GetCommonResources())
	gcpDeploy := gcp.DeploymentProperties{}
	commonUtil.LoadYamlFile(tempDir+"properties.yaml", &gcpDeploy)
	properties := &(gcpDeploy.Resources[0].Properties)
	properties.Name = clusterName
	properties.Zone = zone
	propertiesYaml, _ := yaml.Marshal(gcpDeploy)
	_ = ioutil.WriteFile(tempDir+"properties.yaml", []byte(propertiesYaml), 0644)
}

func setupK8s() {
	logger.Always("Setup Kubernetes Permissions")

	commonUtil.Shell("kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$(gcloud config get-value account)")
	commonUtil.Shell("gcloud services enable cloudbuild.googleapis.com")
	commonUtil.Shell("kubectl create clusterrolebinding permissive-binding --clusterrole=cluster-admin --user=admin --user=kubelet --group=system:serviceaccounts;")
}

func gcpInit(zone string, projectName string) {
	logger.Success("Start GCP init")
	LoginByServiceAccount(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	commonUtil.Shell("gcloud config set compute/zone " + zone)
	commonUtil.Shell("gcloud config set project " + projectName)
	commonUtil.Shell("gcloud services enable compute.googleapis.com")
	commonUtil.Shell("gcloud services enable container.googleapis.com")
}
