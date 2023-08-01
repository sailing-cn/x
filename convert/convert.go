package convert

import "encoding/json"

func ToString(source interface{}) string {
	if source == nil {
		return ""
	}
	byteData, err := json.Marshal(source)
	if err != nil {
		return ""
	}
	return string(byteData)
}
