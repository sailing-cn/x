package iot

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// TimestampString 时间戳：为设备连接平台时的UTC时间，格式为YYYYMMDDHH，如UTC 时间2018/7/24 17:56:20 则应表示为2018072417。
func TimestampString() string {
	strFormatTime := time.Now().Format("2006-01-02 15:04:05")
	strFormatTime = strings.ReplaceAll(strFormatTime, "-", "")
	strFormatTime = strings.ReplaceAll(strFormatTime, " ", "")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func FormatTopic(topic, deviceId string) string {
	return strings.ReplaceAll(topic, "{device_id}", deviceId)
}

// FormatTopic1 下发的时候追加了requestId
func FormatTopic1(topic, deviceId string, requestId string) string {
	str := strings.ReplaceAll(topic, "{device_id}", deviceId)

	return strings.ReplaceAll(str, "#", "request_id ="+requestId)
}

// GetEventTimeStamp 设备采集数据UTC时间（格式：yyyyMMdd'T'HHmmss'Z'），如：20161219T114920Z。
// 设备上报数据不带该参数或参数格式错误时，则数据上报时间以平台时间为准。
func GetEventTimeStamp() int64 {
	now := time.Now().UnixMilli() / 1e6
	return now
}

func Interface2JsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	byteData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(byteData)
}

func GetTopicRequestId(topic string) string {
	return strings.Split(topic, "=")[1]
}

// SmartFileName 根据当前运行的操作系统重新修改文件路径以适配操作系统
func SmartFileName(filename string) string {
	// Windows操作系统适配
	if strings.Contains(OsName(), "windows") {
		pathParts := strings.Split(filename, "/")
		// todo windows 报错
		if len(pathParts) > 1 {
			pathParts[0] = pathParts[0] + ":"
		}
		//pathParts[0] = pathParts[0] + ":"
		return strings.Join(pathParts, "\\")
	}

	return filename
}

func OsName() string {
	return runtime.GOOS
}

func SdkInfo() map[string]string {
	f, err := os.Open("sdk_info")
	if err != nil {
		log.Warning("read sdk info failed")
		return map[string]string{}
	}

	// 文件很小
	info := make(map[string]string)
	buf := bufio.NewReader(f)
	for {
		b, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			log.Warningf("read sdk info failed or end")
			break
		}
		line := string(b)
		if len(line) != 0 {
			parts := strings.Split(line, "=")
			info[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	return info
}

func GetTopicDeviceId(topic string) string {
	//$oc/devices/1458634485842055168_1/sys/properties/report
	return strings.Split(topic, "/")[2]
}
