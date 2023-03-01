package convert

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
)

func ToBytes(in interface{}) (buf []byte) {
	buf, err := json.Marshal(in)
	if err != nil {
		log.Errorf("数据转换失败:%s", err.Error())
		return []byte{}
	}
	return buf
}

func ToInterface(bytes []byte) interface{} {
	var result interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.Errorf("数据转换失败:%s", err.Error())
		return nil
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

// Bytes2Uint16s bytes convert to uint16 for register.
func Bytes2Uint16s(buf []byte) []uint16 {
	data := make([]uint16, 0, len(buf)/2)
	for i := 0; i < len(buf)/2; i++ {
		data = append(data, binary.BigEndian.Uint16(buf[i*2:]))
	}
	return data
}

func Bytes2Uint16(buf []byte) uint16 {
	data := binary.BigEndian.Uint16(buf)
	return data
}

func BytesToDecimal(buf []byte, round int32, scale float32) decimal.Decimal {
	_value := float32(Bytes2Uint16(buf)) * scale
	return decimal.NewFromFloat32(_value).Round(round)
	//result, _ := decimal.NewFromFloat32(_value).Round(round).Float64()
	//return strconv.FormatFloat(result, 'f', 2, 64)
}

func ToIEEE754Float(data []uint16) ([]float32, error) {
	hexs := make([]string, 0)
	for _, item := range data {
		hexs = append(hexs, IntToHex(int64(item)))
	}
	result := make([]float32, 0)
	if len(hexs)%2 != 0 {
		return nil, nil
	}
	for i := 0; i < len(hexs); i += 2 {
		list := hexs[i : i+2]
		str := strings.Join(list, " ")
		n, _ := strconv.ParseUint(str, 16, 32)

		result = append(result, math.Float32frombits(uint32(n)))
	}
	return result, nil

}

func IntToHex(value int64) string {
	return fmt.Sprintf(strconv.FormatInt(value, 16))
}
