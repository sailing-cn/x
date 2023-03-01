package conf_d

type DB struct {
	Connection string `json:"connection"`
	Type       string `json:"type"`
	Debug      bool   `json:"debug"`
}
type Config struct {
	Db *DB `json:"db"`
}
