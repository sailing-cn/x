package conf

type Config struct {
	Db struct {
		Connection string `json:"connection"`
		Type       string `json:"type"`
		Debug      bool   `json:"debug"`
	} `json:"db"`
}
