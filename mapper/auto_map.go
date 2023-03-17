package mapper

import (
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
)

func MapTo(input interface{}, out interface{}) error {
	metadata := &mapstructure.Metadata{}
	err := mapstructure.DecodeMetadata(input, out, metadata)
	//fmt.Printf("keys:%#v unused:%#v\n", metadata.Keys, metadata.Unused)
	return err
}

// MapToProtoMessage 映射到Message
func MapToProtoMessage(input interface{}, out proto.Message) error {
	decoder, err := mapstructure.NewDecoder(protoConfig(out))
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func protoHook(f, t reflect.Kind, data interface{}) (interface{}, error) {
	switch f.String() {
	case "string":
		{
			if t.String() == "ptr" {
				return wrapperspb.String(data.(string)), nil
			}
			break
		}
	case "int":
		{
			if t.String() == "ptr" {
				return wrapperspb.Int32(int32(data.(int))), nil
			}
			break
		}
	case "int64":
		{
			if t.String() == "ptr" {
				return wrapperspb.Int64(data.(int64)), nil
			}
			break
		}
	case "bool":
		{
			if t.String() == "ptr" {
				return wrapperspb.Bool(data.(bool)), nil
			}
			break
		}
	case "float64":
		{
			if t.String() == "ptr" {
				return wrapperspb.Double(data.(float64)), nil
			}
			break
		}
	case "float32":
		{
			if t.String() == "ptr" {
				return wrapperspb.Float(data.(float32)), nil
			}
			break
		}
	case "slice":
		{
			if t.String() == "ptr" {
				return wrapperspb.Bytes(data.([]byte)), nil
			}
			break
		}
	case "uint32":
		{
			if t.String() == "ptr" {
				return wrapperspb.UInt32(data.(uint32)), nil
			}
			break
		}
	case "uint64":
		{
			if t.String() == "ptr" {
				return wrapperspb.UInt64(data.(uint64)), nil
			}
			break
		}
	}
	println("数据类型:%s", f.String())
	return data, nil
}
func protoConfig(d interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		DecodeHook:       protoHook,
		Result:           d,
		WeaklyTypedInput: true,
	}
}
