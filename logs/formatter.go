package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type LogFormatter struct {
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	var data []byte
	if len(entry.Data) > 0 {
		data, _ = json.Marshal(entry.Data)
	}
	newLog = fmt.Sprintf("[%s] [%s] %s %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, data)

	b.WriteString(newLog)
	return b.Bytes(), nil
}
