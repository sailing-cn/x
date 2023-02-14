package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"x/wrong"
)

// Timestamp 时间戳
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// TimeString 时间戳：UTC时间","格式为YYYYMMDDHH","如UTC 时间2018/7/24 17:56:20 则应表示为2018072417。
func TimeString() string {
	strFormatTime := time.Now().Format("2006-01-02 15:04:05")
	strFormatTime = strings.ReplaceAll(strFormatTime, "-", "")
	strFormatTime = strings.ReplaceAll(strFormatTime, " ", "")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

func TimestampString() string {
	timestamp := time.Now().UnixNano() / 1e6
	return strconv.FormatUint(uint64(timestamp), 10)
}

// GetExecPath 获取执行路径
func GetExecPath() string {
	path, err := os.Executable()
	if err != nil {
		log.Printf(err.Error())
	}
	dir := filepath.Dir(path)
	return dir
}

// Interface2JsonString interface转jsonstring
func Interface2JsonString(source interface{}) string {
	if source == nil {
		return ""
	}
	byteData, err := json.Marshal(source)
	if err != nil {
		return ""
	}
	return string(byteData)
}

// ToMap 结构体转map
func ToMap(source interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	buf, _ := json.Marshal(source)
	json.Unmarshal(buf, &data)

	return data
}

func ToMapStr(source interface{}) map[string]string {
	var data = make(map[string]string)
	buf, _ := json.Marshal(source)
	json.Unmarshal(buf, &data)
	return data
}

func ToAny(data interface{}) *anypb.Any {
	gob.Register(data)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(data)
	return &anypb.Any{Value: buf.Bytes()}
}

// PathExists /*判断文件夹是否存在*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathExistsOrCreate(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func LinuxVersion() string {
	cmd := exec.Command("uname", "-a")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := io.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	getVersion01 := fmt.Sprintf("%s", string(opBytes))
	getVersion := strings.Split(getVersion01, " ")
	return getVersion[2]
}

func MD5_SALT(source string, salt string) string {
	b := []byte(source)
	s := []byte(salt)
	h := md5.New()
	h.Write(s) // 先写盐值
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func Name() string {
	p, _ := process.NewProcess(int32(os.Getpid()))
	name, _ := p.Name()
	return name
}

// GetIpv4List 获取本机IPv4 地址列表
func GetIpv4List() ([]string, error) {
	addr, err := net.InterfaceAddrs() //局域文件传输代码
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	result := make([]string, 0)
	for _, address := range addr {
		if ipAddr, ok := address.(*net.IPNet); ok && !ipAddr.IP.IsLoopback() {
			if ipAddr.IP.To4() != nil {
				ipv4 := ipAddr.IP.String()
				result = append(result, ipv4)
			}
		}
	}
	return result, nil
}

func Welcome() {
	//	str := `
	//            DOZ$$$ZOD
	//         NO$777$$$$$7$8N
	//       DZ$7$$$$$77II????78
	//    NO$7$$$$$77$Z$$7?+++++IZN
	//  DO$7$$$$$7$ZD      87++++++78
	// O77$$$$$7$ON          NZI+??78
	//87$$$$77$8       D8888D
	//N$7$$$ZD      N8ZZZZZZZZOD
	//  DZ7Z      DZ$$ZOZZZO8DDOO8N
	//    NO$ODDO$77$$$$ZON     NOZ8N
	//       8Z7777777$ON      DOOZZZD
	//         N8ZZZOD       DOZZZZOZ8
	//  DO$$$ZD           N8ZZZZZZZZZD
	//  DZ$7$$7$ON      NOZZZZZZZZZ8N
	//    NO$7$$77$O88OOZZZZZZZZOD
	//       8Z$$$$$$ZZZZZZZZZOD
	//         N8ZZZZZZZZZZO8N
	//            DOOZZZO8D
	//                                `
	str := `
      ___           ___                                                 ___           ___     
     /  /\         /  /\        ___                       ___          /__/\         /  /\    
    /  /:/_       /  /::\      /  /\                     /  /\         \  \:\       /  /:/_   
   /  /:/ /\     /  /:/\:\    /  /:/      ___     ___   /  /:/          \  \:\     /  /:/ /\  
  /  /:/ /::\   /  /:/~/::\  /__/::\     /__/\   /  /\ /__/::\      _____\__\:\   /  /:/_/::\ 
 /__/:/ /:/\:\ /__/:/ /:/\:\ \__\/\:\__  \  \:\ /  /:/ \__\/\:\__  /__/::::::::\ /__/:/__\/\:\
 \  \:\/:/~/:/ \  \:\/:/__\/    \  \:\/\  \  \:\  /:/     \  \:\/\ \  \:\~~\~~\/ \  \:\ /~~/:/
  \  \::/ /:/   \  \::/          \__\::/   \  \:\/:/       \__\::/  \  \:\  ~~~   \  \:\  /:/ 
   \__\/ /:/     \  \:\          /__/:/     \  \::/        /__/:/    \  \:\        \  \:\/:/  
     /__/:/       \  \:\         \__\/       \__\/         \__\/      \  \:\        \  \::/   
     \__\/         \__\/                                               \__\/         \__\/    
`
	fmt.Println(str)
}

func MysqlDbTypeToCsharp(_type string) string {
	var result string
	switch _type {
	case "bigint":
		result = "long"
		break
	case "mediumint", "int":
		result = "int"
		break
	case "smallint":
		result = "short"
		break
	case "tinyint", "unsigned":
		result = "uint"
		break
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		result = "string"
		break
	case "bit", "bool", "boolean":
		result = "bool"
		break
	case "date", "time", "datetime ", "timestamp", "year":
		result = "DateTime"
		break
	case "decimal", "numeric", "smallmoney", "money":
		result = "decimal"
		break
	case "float":
		result = "float"
		break
	case "image", "binary", "tinyblob", "blob", "mediumblob", "longblob", "varbinary":
		result = "byte[]"
		break
	}
	return result
}

func MysqlDbTypeToGolang(_type string) string {
	switch _type {
	case "tinyint", "bool", "boolean":
		return "bool"
	case "varchar":
		return "string"
	case "datetime":
		return "time.Time"
	case "timestamp":
		return "time.Time"
	case "double":
		return "float64"
	case "float":
		return "float64"
	case "decimal":
		return "float64"
	case "int":
		return "int64"
	case "smallint":
		return "int64"
	case "mediumint":
		return "int64"
	case "integer":
		return "int64"
	case "bigint":
		return "int64"
	default:
		return "string"
	}

}

func MysqlDBtypeToTypescript(_type string) string {
	var result string
	switch _type {
	case "bigint", "mediumint", "int", "smallint", "decimal", "numeric", "smallmoney", "money", "float", "tinyint", "unsigned":
		result = "number"
		break
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		result = "string"
		break
	case "bit", "bool", "boolean":
		result = "boolean"
		break
	case "date", "time", "datetime ", "timestamp", "year":
		result = "Date"
		break
	case "image", "binary", "tinyblob", "blob", "mediumblob", "longblob", "varbinary":
		result = "string"
		break
	}
	return result
}

// RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// RemoveAllChildFile 删除文件下的所有东西
func RemoveAllChildFile(path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		if err, ok := wrong.ShouldOk(os.RemoveAll(filepath.Join(path, entry.Name()))); !ok {
			log.Error(err)
		}
	}
	return nil
}
func CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

// WriteFile 将数据写入文件中(飞追加写入)
func WriteFile(path, version string) error {
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(version)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
