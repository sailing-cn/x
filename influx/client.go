package influx

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	path2 "sailing.cn/v2/path"
	"time"
)

type InfluxConfig struct {
	Influx config
}
type config struct {
	Token   string
	Servers string
	Bucket  string
	Org     string
}

func (conf *InfluxConfig) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(path2.GetExecPath(), "conf.d", "conf.yml")
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

var instance InfluxClient

type InfluxClient interface {
	WritePoint(table string, service string, device string, data map[string]interface{}) error
	GetLast(device string) []map[string]interface{}
	GetClient() influxdb2.Client
}

type influxClient struct {
	Client influxdb2.Client
	bucket string
	org    string
}

func (client influxClient) GetClient() influxdb2.Client {
	return client.Client
}

func (client influxClient) WritePoint(table string, service string, device string, data map[string]interface{}) error {
	api := client.Client.WriteAPIBlocking(client.org, client.bucket)
	P := influxdb2.NewPoint(table,
		map[string]string{"service_id": service, "device_id": device},
		data, time.Now())
	return api.WritePoint(context.Background(), P)
}

func (client influxClient) GetLast(device string) []map[string]interface{} {
	api := client.Client.QueryAPI(client.org)
	result, _ := api.Query(context.Background(),
		fmt.Sprintf(`from(bucket:"%s") |> range(start:-30d) 
                    |> filter(fn: (r) => r["device_id"] == "%s")
					|> last()`,
			client.bucket, device))
	list := make(map[string][]interface{})
	datas := make([]interface{}, 0)
	for result.Next() {
		field := result.Record().Field()
		value := result.Record().Value()
		hasKey := false
		for key := range list {
			if field == key {

				hasKey = true
				break
			}
		}
		if hasKey == false {
			datas = append(datas, value)
		}
		list[field] = datas
		datas = make([]interface{}, 0)
	}
	points := make([]map[string]interface{}, 0)
	keys := make([]string, 0)
	for key := range list {
		keys = append(keys, key)
	}

	for key, values := range list {
		index := 0
		for _, value := range values {
			var p map[string]interface{}
			if len(points) > index {
				p = points[index]
			}
			if p == nil {
				p = make(map[string]interface{})
				points = append(points, p)
			}
			p[key] = value
			points[index] = p
			index++
		}
	}
	return points
}

func Get() InfluxClient {
	if instance == nil {
		cnf := &InfluxConfig{}
		cnf.Init("")
		instance = CreateClientWithClient(cnf)
	}
	return instance
}

func CreateClientWithClient(conf *InfluxConfig) influxClient {
	c := influxdb2.NewClient(conf.Influx.Servers, conf.Influx.Token)
	return influxClient{Client: c, org: conf.Influx.Org, bucket: conf.Influx.Bucket}
}
