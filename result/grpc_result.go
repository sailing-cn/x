package result

//
//import (
//	"sailing.cn/convert"
//	"sailing.cn/protobuf/core"
//)
//
//func GrpcSuccess(data interface{}) *core.Result {
//	var result = &core.Result{}
//	result.Code = 1
//	result.Message = "success"
//	result.Success = true
//	result.Data = convert.ToBytes(data)
//	return result
//}
//
//func GrpcFail(message string, data interface{}) *core.Result {
//	var result = &core.Result{}
//	result.Code = 2
//	result.Message = message
//	result.Success = false
//	result.Data = convert.ToBytes(data)
//	return result
//}
