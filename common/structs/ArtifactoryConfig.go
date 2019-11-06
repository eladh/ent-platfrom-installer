package structs

type ArtifactoryConfig struct {
	Ingress struct {
		Annotations struct {
			CertmanagerK8SIoClusterIssuer                string `yaml:"certmanager.k8s.io/cluster-issuer"`
			IngressKubernetesIoForceSslRedirect          string `yaml:"ingress.kubernetes.io/force-ssl-redirect"`
			IngressKubernetesIoProxyBodySize             string `yaml:"ingress.kubernetes.io/proxy-body-size"`
			IngressKubernetesIoProxyReadTimeout          string `yaml:"ingress.kubernetes.io/proxy-read-timeout"`
			IngressKubernetesIoProxySendTimeout          string `yaml:"ingress.kubernetes.io/proxy-send-timeout"`
			KubernetesIoIngressClass                     string `yaml:"kubernetes.io/ingress.class"`
			KubernetesIoTLSAcme                          string `yaml:"kubernetes.io/tls-acme"`
			NginxIngressKubernetesIoConfigurationSnippet string `yaml:"nginx.ingress.kubernetes.io/configuration-snippet"`
			NginxIngressKubernetesIoProxyBodySize        string `yaml:"nginx.ingress.kubernetes.io/proxy-body-size"`
		} `yaml:"annotations"`
		DefaultBackend struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"defaultBackend"`
		Enabled bool     `yaml:"enabled"`
		Hosts   []string `yaml:"hosts"`
		TLS     []struct {
			Hosts      []string `yaml:"hosts"`
			SecretName string   `yaml:"secretName"`
		} `yaml:"tls"`
	} `yaml:"ingress"`
	Nginx struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"nginx"`
	Postgresql struct {
		Enabled     bool `yaml:"enabled"`
		Persistence struct {
			Size string `yaml:"size"`
		} `yaml:"persistence"`
		PostgresPassword string `yaml:"postgresPassword"`
	} `yaml:"postgresql"`
}
