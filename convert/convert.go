package convert

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

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

func ToBytes(in interface{}) (buf []byte) {
	buf, err := json.Marshal(in)
	if err != nil {
		return []byte{}
	}
	return buf
}

// ToTimestamp 将指定格式的时间转换成时间戳 2022-03-11 11-30-50   03091545    0803062323
func ToTimestamp(str string) int64 {
	var result int64 = 0
	if strings.Contains(str, "-") || strings.Contains(str, " ") {
		str = strings.ReplaceAll(str, " ", "-")
	} else {
		if len(str) == 8 {
			str = str + "00"
		}
		for i := 0; i < len(str); i += 3 {
			str = str[:i] + "-" + str[i:]
		}
		str = strconv.Itoa(time.Now().Year()) + str
	}

	strFormatTimeList := strings.Split(str, "-")
	if len(strFormatTimeList) == 6 {
		year, err := strconv.Atoi(strFormatTimeList[0])
		month, err := strconv.Atoi(strFormatTimeList[1])
		day, err := strconv.Atoi(strFormatTimeList[2])
		hour, err := strconv.Atoi(strFormatTimeList[3])
		min, err := strconv.Atoi(strFormatTimeList[4])
		sec, err := strconv.Atoi(strFormatTimeList[5])
		if err != nil {
			return result
		}
		result = time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local).Unix()
	}
	return result
}

func ToModel[T interface{}](bytes []byte) *T {
	var result T
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.Errorf("数据转换失败:%s", err.Error())
		return nil
	}
	return &result
}
