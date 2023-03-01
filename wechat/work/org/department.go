package org

import (
	"fmt"
	"sailing.cn/wechat/util"
)

const (
	FetchDepartmentListURL = "https://qyapi.weixin.qq.com/cgi-bin/department/simplelist"
)

func (c *Client) GetDepartmentList(departmentId string) (interface{}, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	var response []byte
	response, err = util.HTTPGet(fmt.Sprintf("%s?access_token=%v&id=%s", FetchDepartmentListURL, accessToken, departmentId))
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = util.DecodeWithError(response, result, "GetDepartmentList")
	if err != nil {
		return nil, err
	}
	return result, nil
}
