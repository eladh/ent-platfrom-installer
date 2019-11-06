package structs

type AddArtifactoryTrustedKeyPayout struct {
	Alias     string `json:"alias"`
	PublicKey string `json:"public_key"`
}

type Storage struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Type     string `yaml:"type"`
}

type Artifactory struct {
	Name       string           `yaml:"name"`
	Site       string           `yaml:"site"`
	AuthServer bool             `yaml:"auth_server"`
	HA         HighAvailability `yaml:"high_availability"`
	Storage    Storage          `yaml:"storage"`
	Repos 	   []Repo			`yaml:"repos"`
	Replications Replications   `yaml:"replications"`
}

type HighAvailability struct {
	ReplicaCount int `yaml:"replica_count"`
	MinAvailable int `yaml:"min_available"`
}

type ArtifactoryHost struct {
	Name string `json:"name"`
}

type ArtifactoryLicense struct {
	LicenseKey string `json:"licenseKey"`
}

type Repo struct {
	Local bool `yaml:"local"`
	Remote bool `yaml:"remote"`
	Virtual bool `yaml:"virtual"`
	PackageType string `yaml:"package_type"`
	Name string `yaml:"name"`
	Url string `yaml:"url"`
	NpmRemoteUrl string `yaml:"externalDependenciesRemoteRepo"`

}

type Replications []struct {
	Source string `yaml:"source"`
	Repo   string `yaml:"repo"`
}