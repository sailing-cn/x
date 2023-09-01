package conf

type Config struct {
	Db *DbConfig `json:"db"`
}

type DbConfig struct {
	Connection string `json:"connection"`
	Type       string `json:"type"`
	Debug      bool   `json:"debug"`
}
