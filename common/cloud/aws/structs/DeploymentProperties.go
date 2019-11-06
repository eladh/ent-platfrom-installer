package aws

type NodeGroup struct {
	Name         string `yaml:"name"`
	InstanceType string `yaml:"instanceType"`
	MinSize      int    `yaml:"minSize"`
	VolumeSize   int    `yaml:"volumeSize"`
	VolumeType   string `yaml:"volumeType"`
	MaxSize      int    `yaml:"maxSize"`
	AllowSSH     bool   `yaml:"allowSSH"`
	Ami          string `yaml:"ami"`
	Labels       struct {
		NodeGroupType string `yaml:"nodegroup-type"`
	} `yaml:"labels"`
	PreBootstrapCommand []string `yaml:"preBootstrapCommand"`
}

type DeploymentProperties struct {
	ApiVersion        string   `yaml:"apiVersion"`
	Kind              string   `yaml:"kind"`
	AvailabilityZones []string `yaml:"availabilityZones"`
	Metadata          struct {
		Name   string `yaml:"name"`
		Region string `yaml:"region"`
	} `yaml:"metadata"`
	NodeGroups []NodeGroup `yaml:"nodeGroups"`
}
