package gcp

import (
	"common/kubernetes"
	"common/structs"
	commmonUtils "common/utils"
	"github.com/kris-nova/logger"
	"io/ioutil"
	"os"
)

func LoginByServiceAccount(licenseLocation string) {
	commmonUtils.Shell("gcloud auth activate-service-account --key-file=" + licenseLocation)
}

func Connect2Cluster(clusterName string ,project string) {
	commmonUtils.Shell("gcloud container clusters get-credentials " + clusterName + " --zone us-central1-a --project " + project)
}

func CreateInternalDnsRecords(setupInfo structs.SetupInfo , host string) {
	cluster := setupInfo.Cluster.Name
	project := setupInfo.Vendor.Project
	domain := setupInfo.Cluster.Domain

	alias := host + "." + cluster + "." + domain + "."
	logger.Always("create dns for domain " + alias)
	address := kubernetes.GetServiceIp(host + "-artifactory-nginx")
	_ = os.Remove("./transaction.yaml")

	commmonUtils.ShellCurrentDir("gcloud beta dns --project=" + project + " managed-zones create " + cluster +
		" --description=\"new zone\" --dns-name=\"" + cluster + "." + domain + "." + "\" --visibility=\"private\" --networks \"default\"" , setupInfo.TempDir)

	commmonUtils.ShellCurrentDir("gcloud dns --project=" + project + " record-sets transaction start --zone=" + cluster , setupInfo.TempDir)
	commmonUtils.ShellCurrentDir("gcloud dns --project=" + project + " record-sets transaction add " + address + " --name=" + alias + " --ttl=300 --type=A --zone=" + cluster , setupInfo.TempDir)
	commmonUtils.ShellCurrentDir("gcloud dns --project=" + project + " record-sets transaction add " + address + " --name=docker." + alias + " --ttl=300 --type=A --zone=" + cluster , setupInfo.TempDir)
	commmonUtils.ShellCurrentDir("gcloud dns --project=" + project + " record-sets transaction execute --zone=" + cluster , setupInfo.TempDir)
}

func DeployCertToClusterNodes(zone string, host string, tempDir string) {
	certificate, _ := ioutil.ReadFile(tempDir + "ca.crt")
	nodes, _ := kubernetes.GetNodes()

	for _, node := range nodes.Items {
		command := "gcloud compute ssh --zone \"" + zone + "\" " + node.Name +
			" --command \" sudo su -c ' mkdir -p /etc/docker/certs.d/" + host + "' root &&" +
			" echo '" + string(certificate) + "' > /tmp/ca.crt && sudo su -c ' mv /tmp/ca.crt /etc/docker/certs.d/" +
			host + "/' root\"  </dev/null;"

		commmonUtils.Shell(command)
	}
}


func GetAddressByNameAndRegion(name string ,region string) string {
	return commmonUtils.Shell("gcloud compute addresses describe " +  name +
		" --region " + region + " | head -1 | cut -d \":\" -f2 |  tr -d '[:space:]'")
}

func CreateExternalIPAddress(serviceName string ,setupInfo structs.SetupInfo) string  {
	commmonUtils.Shell("gcloud compute addresses create "+  serviceName + " --project=" + setupInfo.Vendor.Project +
		" --region=" + setupInfo.Vendor.Region)
	return GetAddressByNameAndRegion(serviceName ,setupInfo.Vendor.Region)
}
