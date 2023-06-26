package id

import (
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Error(err)
		return
	}
}

func SnowflakeString() string {
	return node.Generate().String()
}

func Snowflake() int64 {
	return node.Generate().Int64()
}
