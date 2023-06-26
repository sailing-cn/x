package timestamp

import (
	"strconv"
	"time"
)

// Timestamp 时间戳
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// TimestampString 时间戳字符串
func TimestampString() string {
	timestamp := Timestamp()
	return strconv.FormatUint(uint64(timestamp), 10)
}
