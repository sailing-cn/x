package client

import (
	"errors"
	"fmt"
	"net/http"
	sdk "sailing.cn/emqx/v5/http"
	"sailing.cn/utils"
)

const (
	CLIENT_URL = "/clients"
)

type ClientService sdk.Service

func (c *ClientService) ListClient(query *ClientPageQuery) (*sdk.PageResponse[ClientResponse], error) {
	result := new(sdk.PageResponse[ClientResponse])
	_query := utils.ToMapStr(*query)
	resp, err := c.Client.R().
		SetQueryParams(_query).
		SetResult(&result).Get(c.Client.RequestURL(CLIENT_URL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	return result, nil
}
