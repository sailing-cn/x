package mapper

import (
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"reflect"
)

func MapTo(input interface{}, out interface{}) error {
	metadata := &mapstructure.Metadata{}
	err := mapstructure.DecodeMetadata(input, out, metadata)
	//fmt.Printf("keys:%#v unused:%#v\n", metadata.Keys, metadata.Unused)
	return err
}

// ProtoToStruct proto转结构体
func ProtoToStruct(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: func(f, t reflect.Type, data interface{}) (interface{}, error) {
			from := f.String()
			to := t.String()
			if from == "*wrapperspb.StringValue" && to == "string" {
				source := data.(*wrapperspb.StringValue)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.Int64Value" && to == "int64" {
				source := data.(*wrapperspb.Int64Value)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.Int32Value" && (to == "int" || to == "int32") {
				source := data.(*wrapperspb.Int32Value)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.UInt32Value" && to == "uint32" {
				source := data.(*wrapperspb.UInt32Value)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.UInt64Value" && to == "uint64" {
				source := data.(*wrapperspb.UInt64Value)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.BytesValue" && to == "[]byte" {
				source := data.(*wrapperspb.BytesValue)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.DoubleValue" && to == "float64" {
				source := data.(*wrapperspb.DoubleValue)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.FloatValue" && to == "float32" {
				source := data.(*wrapperspb.FloatValue)
				return source.GetValue(), nil
			} else if from == "*wrapperspb.BoolValue" && to == "bool" {
				source := data.(*wrapperspb.BoolValue)
				return source.GetValue(), nil
			}

			return data, nil
		},
		Result:           out,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// MapToProtoMessage 映射到Message
func MapToProtoMessage(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook:       protoDecoderHook,
		Result:           out,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func MapToProtoMessageList[T interface{}](input []interface{}, out []interface{}) error {
	if input == nil {
		return nil
	}
	for _, item := range input {
		var m T
		MapToProtoMessage(item, &m)
		out = append(out)
	}
	return nil
}

func protoDecoderHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	from := f.String()
	to := t.String()
	if data == nil {
		return nil, nil
	}
	if from == "string" && to == "*wrapperspb.StringValue" {
		return &wrapperspb.StringValue{Value: data.(string)}, nil
	} else if from == "int64" && to == "*wrapperspb.Int64Value" {
		source := data.(int64)
		return source, nil
	} else if (from == "int" || from == "int32") && to == "*wrapperspb.Int32Value" {
		source := data.(int)
		return source, nil
	} else if from == "uint32" && to == "*wrapperspb.UInt32Value" {
		source := data.(uint32)
		return source, nil
	} else if from == "uint64" && to == "*wrapperspb.UInt64Value" {
		source := data.(uint64)
		return source, nil
	} else if from == "[]byte" && to == "*wrapperspb.BytesValue" {
		source := data.([]byte)
		return source, nil
	} else if from == "float64" && to == "*wrapperspb.DoubleValue" {
		source := data.(float64)
		return source, nil
	} else if from == "float32" && to == "*wrapperspb.FloatValue" {
		source := data.(float32)
		return source, nil
	} else if from == "bool" && to == "*wrapperspb.BoolValue" {
		source := data.(bool)
		return source, nil
	}
	return data, nil
}
