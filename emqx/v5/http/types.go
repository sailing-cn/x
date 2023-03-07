package http

type PageResponse[T interface{}] struct {
	Data []T  `json:"data"`
	Page Meta `json:"meta"`
}

type Meta struct {
	Count int `json:"count"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

// Metrics 监测
type Metrics struct {
	Failed     int     `json:"failed"`
	Matched    int     `json:"matched"`
	Rate       float32 `json:"rate"`
	RateLast5M float32 `json:"rate_last5m"`
	RateMax    float32 `json:"rate_max"`
	Success    int     `json:"success"`
}
