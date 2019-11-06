package gcp

type DeploymentProperties struct {
	Imports []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
	} `yaml:"imports"`
	Resources []struct {
		Name       string `yaml:"name"`
		Type       string `yaml:"type"`
		Properties struct {
			Name             string `yaml:"name"`
			Description      string `yaml:"description"`
			Zone             string `yaml:"zone"`
			InitialNodeCount int    `yaml:"initialNodeCount"`
			MinNodeCount     int    `yaml:"minNodeCount"`
			MaxNodeCount     int    `yaml:"maxNodeCount"`
			MachineType      string `yaml:"machineType"`
		} `yaml:"properties"`
	} `yaml:"resources"`
}
