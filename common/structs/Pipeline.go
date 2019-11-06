package structs

import "time"


type PipelineSource struct {
	ProjectIntegrationID int    `json:"projectIntegrationId"`
	RepositoryFullName   string `json:"repositoryFullName"`
	ProjectID            int    `json:"projectId"`
	Branch               string `json:"branch"`
	FileFilter           string `json:"fileFilter"`
}

type PipelineChart struct {
	API struct {
		Artifactory struct {
			Password string      `yaml:"password"`
			URL      interface{} `yaml:"url"`
			Username string      `yaml:"username"`
		} `yaml:"artifactory"`
		ExternalURL interface{} `yaml:"externalUrl"`
		RabbitMq    struct {
			UIUserPassword string `yaml:"uiUserPassword"`
			UserPassword   string `yaml:"userPassword"`
		} `yaml:"rabbitMq"`
		Service struct {
			LoadBalancerIP interface{} `yaml:"loadBalancerIP"`
			Type           string      `yaml:"type"`
		} `yaml:"service"`
		ServiceID string `yaml:"serviceId"`
		Token     string `yaml:"token"`
	} `yaml:"api"`
	Postgresql struct {
		Enabled            bool   `yaml:"enabled"`
		PostgresqlDatabase string `yaml:"postgresqlDatabase"`
		PostgresqlPassword string `yaml:"postgresqlPassword"`
		PostgresqlUsername string `yaml:"postgresqlUsername"`
	} `yaml:"postgresql"`
	Nodes struct {
		QuotaSize            int   `yaml:"quotaSize"`
		PoolSize 			 int   `yaml:"poolSize"`
		ReplicaCount         int   `yaml:"replicaCount"`
	} `yaml:"nodes"`
	RunMode string `yaml:"runMode"`
	Www     struct {
		ExternalURL interface{} `yaml:"externalUrl"`
		Service     struct {
			LoadBalancerIP interface{} `yaml:"loadBalancerIP"`
			Port           int         `yaml:"port"`
			Type           string      `yaml:"type"`
		} `yaml:"service"`
	} `yaml:"www"`
}


type PipelineIntegration struct {
	Name        string `json:"name"`
	ProjectID   int    `json:"projectId"`
	PropertyBag struct {
	} `json:"propertyBag"`
	MasterIntegrationID   int              `json:"masterIntegrationId"`
	MasterIntegrationName string           `json:"masterIntegrationName"`
	MasterIntegrationType string           `json:"masterIntegrationType"`
	FormJSONValues        []LabelValuePair `json:"formJSONValues"`
}


type LabelValuePair struct {
	Label string `json:"label"`
	Value string `json:"value"`
}


type PipelineIntegrationResponse struct {
	PropertyBag struct {
	} `json:"propertyBag"`
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	ProjectID             int       `json:"projectId"`
	MasterIntegrationName string    `json:"masterIntegrationName"`
	MasterIntegrationType string    `json:"masterIntegrationType"`
	MasterIntegrationID   int       `json:"masterIntegrationId"`
	CreatedByUserName     string    `json:"createdByUserName"`
	UpdatedByUserName     string    `json:"updatedByUserName"`
	CreatedBy             int       `json:"createdBy"`
	UpdatedBy             int       `json:"updatedBy"`
	ProviderID            int       `json:"providerId"`
	UpdatedAt             time.Time `json:"updatedAt"`
	CreatedAt             time.Time `json:"createdAt"`
	FormJSONValues        []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"formJSONValues"`
}