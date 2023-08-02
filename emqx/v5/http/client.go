package http

import (
	"github.com/imroc/req/v3"
	"net/url"
	"sailing.cn/emqx/v5/config"
)

type APIClient struct {
	cfg      *config.Configuration
	BasePath *url.URL
	client   *req.Client
}

type Service struct {
	Client *APIClient
}

func NewAPIClient(cfg *config.Configuration) *APIClient {
	c := &APIClient{}
	c.cfg = cfg
	c.BasePath, _ = url.Parse(cfg.Server + "/api/" + cfg.ApiVersion)
	c.client = req.C()
	c.client.SetBaseURL(c.BasePath.String())
	return c
}

func (c *APIClient) RequestURL(route string) string {
	u := *c.BasePath
	u.Path = u.Path + route
	return u.String()
}

func (c *APIClient) R() *req.Request {
	return c.client.R().SetBasicAuth(c.cfg.ApiKey, c.cfg.ApiSecret)
}
