package mapper

import (
	"github.com/mitchellh/mapstructure"
)

func MapTo(input interface{}, out interface{}) error {
	metadata := &mapstructure.Metadata{}
	err := mapstructure.DecodeMetadata(input, out, metadata)
	//fmt.Printf("keys:%#v unused:%#v\n", metadata.Keys, metadata.Unused)
	return err
}
