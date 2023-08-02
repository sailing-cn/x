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

func ToMapString(source interface{}) map[string]string {
	var data = make(map[string]string)
	buf, _ := json.Marshal(source)
	_ = json.Unmarshal(buf, &data)
	return data
}
