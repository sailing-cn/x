package context

import (
	"github.com/imroc/req/v3"
	"net/url"
	"sailing.cn/zentao/config"
	"strings"
)

type Context struct {
	*config.Config
	Token   string
	baseURL *url.URL
	Client  *req.Client
}

func (c *Context) SetBaseURL(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseURL = baseURL

	return nil
}

func (c *Context) RequestURL(path string) string {
	u := *c.baseURL
	u.Path = c.baseURL.Path + path
	return u.String()
}

func (c *Context) SetClient() {
	c.Client = req.C()
}
