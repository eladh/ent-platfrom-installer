package aws

import (
	"common/cloud/aws/structs"
	"common/structs"
	"common/utils"
	"github.com/kris-nova/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Install(setupInfo structs.SetupInfo) {
	zone := setupInfo.Vendor.Zone
	clusterName := setupInfo.Cluster.Name

	installCluster(clusterName, zone, setupInfo.TempDir)
}

func Uninstall(setupInfo structs.SetupInfo) {
	zone := setupInfo.Vendor.Zone
	clusterName := setupInfo.Cluster.Name
	utils.Shell("eksctl delete cluster --name=" + clusterName + "  --region=" + zone)
}

func installCluster(clusterName string, region string, tempDir string) {
	logger.Always("Install AWS EKS Cluster")

	utils.CopyResourceToLocation("cloud/aws/group-nodes.yaml" ,tempDir+"group-nodes.yaml" ,utils.GetCommonResources())

	awsDeploy := aws.DeploymentProperties{}
	utils.LoadYamlFile(tempDir+"./group-nodes.yaml", &awsDeploy)

	awsDeploy.Metadata.Name = clusterName + "-cluster"
	awsDeploy.Metadata.Region = region

	nodeGroup := &(awsDeploy.NodeGroups[0])
	nodeGroup.Name = clusterName + "-public"
	nodeGroup.Labels.NodeGroupType = clusterName + "-workloads"

	propertiesYaml, _ := yaml.Marshal(awsDeploy)

	_ = ioutil.WriteFile(tempDir+"/group-nodes.yaml", []byte(propertiesYaml), 0644)
	utils.Shell("eksctl create cluster --config-file " + tempDir + "/group-nodes.yaml")
}
