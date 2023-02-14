package wrong

import "fmt"

var (
	DATA_NOT_FOUND = "该数据不存在或已删除"
	DATA_EXIST     = "已经存在"
	QUERY_ERROR    = "查询异常"
)

func DataExist(src string) string {
	return fmt.Sprintf("%s%s", src, DATA_EXIST)
}
