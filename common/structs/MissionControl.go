package structs

type BaseUrlPayload struct {
	BaseURL string `json:"base_url"`
}

type UpdateJfmcScriptsPayload struct {
	URL           string `json:"url"`
	Branch        string `json:"branch"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Enabled       bool   `json:"enabled"`
	Bidirectional bool   `json:"bidirectional"`
	Restore       bool   `json:"restore"`
}

type SitePayload struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	City        City   `yaml:"city" json:"city"`
}

type City struct {
	Name        string  `yaml:"name" json:"name"`
	CountryCode string  `yaml:"country_code" json:"country_code"`
	Latitude    float64 `yaml:"latitude" json:"latitude"`
	Longitude   float64 `yaml:"longitude" json:"longitude"`
}

type AttachLicencePayload struct {
	ServiceName string `json:"service_name"`
	Deploy      bool   `json:"deploy"`
	Instances   int    `json:"number_of_licenses"`
}

type AddService struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	URL                  string `json:"url"`
	SiteName             string `json:"site_name"`
	UserName             string `json:"username"`
	Password             string `json:"password"`
	PairWithAuthProvider bool   `json:"pair_with_auth_provider"`
	AuthProvider         bool   `json:"auth_provider"`
	Type                 string `json:"type"`
}

type AddAccessFederationPayload struct {
	Entities []string          `json:"entities"`
	Targets  []ArtifactoryHost `json:"targets"`
}

type AddRepoPayload struct {
	Type  string  `json:"rclass"`
	PackageType  string  `json:"packageType"`
	URL  string  `json:"url"`
	Repositories  []string  `json:"repositories"`
	ExternalDependencyEnabled  bool  `json:"externalDependenciesEnabled"`
	DefaultDeploymentRepo  string  `json:"defaultDeploymentRepo"`
}

type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}


type AddReplicationPayload struct {
}
