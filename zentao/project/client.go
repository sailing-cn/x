package project

import "zentao/context"

type Client struct {
	ctx *context.Context
}

func NewClient(ctx *context.Context) *Client {
	return &Client{ctx: ctx}
}
