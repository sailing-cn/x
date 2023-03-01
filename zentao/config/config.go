package config

import "sailing.cn/zentao/cache"

type Config struct {
	Domain   string `json:"domain" yaml:"domain"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Cache    cache.Cache
}
