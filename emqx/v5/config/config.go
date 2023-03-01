package config

type Configuration struct {
	Server     string `json:"server,omitempty"`
	ApiKey     string `json:"key,omitempty"`
	ApiSecret  string `json:"secret,omitempty"`
	ApiVersion string `json:"version,omitempty"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (cnf *Configuration) Init() error {
	return nil
}
