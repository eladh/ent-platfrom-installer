package structs

type SonarConfig struct {
	SonarQubeServers struct {
		Installations struct {
			MySonarQube struct {
				MojoVersion               float64     `yaml:"mojoVersion"`
				ServerAuthenticationToken interface{} `yaml:"serverAuthenticationToken"`
				ServerURL                 interface{} `yaml:"serverUrl"`
				ServerVersion             string      `yaml:"serverVersion"`
				Triggers                  struct {
					SkipScmCause      bool `yaml:"skipScmCause"`
					SkipUpstreamCause bool `yaml:"skipUpstreamCause"`
				} `yaml:"triggers"`
			} `yaml:"my-sonar-qube"`
		} `yaml:"installations"`
	} `yaml:"sonar_qube_servers"`
}
