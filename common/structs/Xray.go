package structs

type Watch struct {
	GeneralData      GeneralData      `json:"general_data" yaml:"general_data"`
	ProjectResources ProjectResources `json:"project_resources" yaml:"project_resources"`
	AssignedPolicies []AssignedPolicy `json:"assigned_policies" yaml:"assigned_policies"`
}

type GeneralData struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Active      bool   `json:"active" yaml:"active"`
}
type Filters struct {
	Type  string `json:"type" yaml:"type"`
	Value string `json:"value" yaml:"value"`
}

type Resources struct {
	Type      string    `json:"type" yaml:"type" `
	BinMgrID  string    `json:"bin_mgr_id" yaml:"bin_mgr_id"`
	Name      string    `json:"name" yaml:"name"`
	Filters   []Filters `json:"filters,omitempty" yaml:"filters,omitempty"`
	Clickable bool      `json:"clickable,omitempty" yaml:"clickable,omitempty"`
}

type ProjectResources struct {
	Resources []Resources `json:"resources" yaml:"resources"`
}

type AssignedPolicy struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

type Policy struct {
	Name        string  `json:"name" yaml:"name"`
	Type        string  `json:"type" yaml:"type"`
	Description string  `json:"description" yaml:"description"`
	Rules       []Rules `json:"rules" yaml:"rules"`
}

type Criteria struct {
	MinSeverity string `json:"min_severity" yaml:"min_severity"`
}

type BlockDownload struct {
	Unscanned bool `json:"unscanned" yaml:"unscanned"`
	Active    bool `json:"active" yaml:"active"`
}

type Actions struct {
	FailBuild     bool          `json:"fail_build" yaml:"fail_build"`
	BlockDownload BlockDownload `json:"block_download" yaml:"block_download"`
}

type Rules struct {
	Name     string   `json:"name" yaml:"name"`
	Priority int      `json:"priority" yaml:"priority"`
	Criteria Criteria `json:"criteria" yaml:"criteria"`
	Actions  Actions  `json:"actions" yaml:"actions"`
}

type IndexedArtifactoryPayload struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	BinMgrURL    string `json:"binMgrUrl"`
	BinMgrID     string `json:"binMgrId"`
	BinMgrDesc   string `json:"binMgrDesc"`
	ProxyEnabled bool   `json:"proxy_enabled"`
}

type WizardPayload struct {
	BaseURL  string `json:"base_url"`
	Finished bool   `json:"finished"`
	StepNum  int    `json:"step_num"`
}

type IndexedRepo struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	PackageType string `yaml:"pkg_type"`
}

type LoginPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type XrayServer struct {
	Name        string         `yaml:"name"`
	Site        string         `yaml:"site"`
	Artifactory string         `yaml:"artifactory"`
	Builds      []string       `yaml:"builds"`
	Repos       []IndexedRepo  `yaml:"repos"`
	Watches		[]Watch		   `yaml:"watches"`
	Policies    []Policy	   `yaml:"policies"`
}

type UpdateIndexedRepos struct {
	Repos []IndexedRepo `json:"indexed_repos"`
}

type UpdateIndexedBuilds struct {
	Builds []string `json:"indexed_builds"`
}
