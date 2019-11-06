package structs

type Param struct {
	Name         string `yaml:"name" json:"name"`
	Type         string `yaml:"type" json:"type"`
	DefaultValue string `yaml:"default_value" json:"default_value"`
	Desc         string `yaml:"desc" json:"desc"`
}

type Job struct {
	Name     string  `yaml:"name" json:"name"`
	Url      string  `yaml:"url" json:"url"`
	Pipeline string  `yaml:"pipeline" json:"pipeline"`
	Params   []Param `yaml:"params" json:"params"`
}

type SetupInfo struct {
	TempDir string
	Cluster struct {
		Name   string `yaml:"name"`
		Domain string `yaml:"domain"`
	} `yaml:"cluster"`
	Pipeline struct {
		GitToken  string `yaml:"token"`
		GitRepo   string `yaml:"repo"`
		GitSource string `yaml:"source"`
	} `yaml:"pipeline"`
	Vendor struct {
		Type    string `yaml:"type"`
		Zone    string `yaml:"zone"`
		Region  string `yaml:"region"`
		Project string `yaml:"project"`
		Gcp     struct {
			Storage struct {
				Identity string `yaml:"identity"`
				Secret   string `yaml:"secret"`
			} `yaml:"storage"`
		} `yaml:"gcp"`
	} `yaml:"vendor"`
	ArtifactoryLicense string        `yaml:"art_license"`
	EdgeLicense        string        `yaml:"edge_license"`
	Sites              []SitePayload `yaml:"sites"`
	Services           struct {
		InstallOnly  bool          `yaml:"install_only"`
		Artifactory  []Artifactory `yaml:"artifactory"`
		Edges        []Artifactory `yaml:"edges"`
		Xray         []XrayServer  `yaml:"xray"`
		Distribution struct {
			Name string `yaml:"name"`
			Site string `yaml:"site"`
		} `yaml:"distribution"`
		Versions struct {
			Artifactory   string `yaml:"artifactory"`
			ArtifactoryHA string `yaml:"artifactory_ha"`
			Xray          string `yaml:"xray"`
			Distribution  string `yaml:"distribution"`
			Jfmc          string `yaml:"jfmc"`
			Sonar         string `yaml:"sonar"`
			Jenkins       string `yaml:"jenkins"`
		} `yaml:"versions"`
	} `yaml:"services"`
	Tools struct {
		Dev       bool `yaml:"dev"`
		Glowroot  bool `yaml:"glowroot"`
		Sonarqube bool `yaml:"sonarqube"`
		Jenkins   struct {
			Jobs []Job  `yaml:"jobs"`
			Site string `yaml:"site"`
		} `yaml:"jenkins"`
		Ldap struct {
		} `yaml:"ldap"`
	} `yaml:"tools"`
}
