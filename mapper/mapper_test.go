package mapper

import (
	"encoding/json"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sailing.cn/utils"
	"testing"
)

type Person struct {
	Name         string
	Age          int
	Birthday     string
	Sex          bool
	CreationTime int64
	Double       float64
	UInt32       uint32
	UInt64       uint64
	Bytes        []byte
	Float        float32
}

type ClientPersonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Age          *wrapperspb.Int32Value  `protobuf:"bytes,1,opt,name=Age,proto3" json:"Age,omitempty"`                   //标识
	Name         *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`                 //名称
	Birthday     string                  `protobuf:"bytes,3,opt,name=Birthday,proto3" json:"Birthday,omitempty"`         //图标
	Sex          *wrapperspb.BoolValue   `protobuf:"bytes,4,opt,name=Sex,proto3" json:"Sex,omitempty"`                   //开发商
	CreationTime *wrapperspb.Int64Value  `protobuf:"bytes,5,opt,name=CreationTime,proto3" json:"CreationTime,omitempty"` //开发商
	Double       *wrapperspb.DoubleValue `protobuf:"bytes,5,opt,name=Double,proto3" json:"Double,omitempty"`
	UInt64       *wrapperspb.UInt64Value `protobuf:"bytes,6,opt,name=UInt64,proto3" json:"UInt64,omitempty"`
	UInt32       *wrapperspb.UInt32Value `protobuf:"bytes,7,opt,name=UInt32,proto3" json:"UInt32,omitempty"`
	Bytes        *wrapperspb.BytesValue  `protobuf:"bytes,8,opt,name=Bytes,proto3" json:"Bytes,omitempty"`
	Float        *wrapperspb.FloatValue  `protobuf:"bytes,9,opt,name=Float,proto3" json:"Float,omitempty"`
}

var file_client_proto_msgTypes = make([]protoimpl.MessageInfo, 7)

func (x *ClientPersonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func TestMapToProto(t *testing.T) {
	person := &Person{
		Name:         "张三",
		Age:          18,
		Birthday:     utils.TimestampString(),
		Sex:          false,
		CreationTime: utils.Timestamp(),
		Float:        1.32,
		UInt64:       64,
		UInt32:       32,
		Bytes:        []byte("kkkkkkkk"),
		Double:       1.64,
	}

	request := &ClientPersonRequest{}
	err := MapToProtoMessage(person, request)
	if err != nil {
		t.Fatal(err)
	} else {
		msg, _ := json.Marshal(request)
		t.Log(string(msg))
	}
}
