package emqx

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"path/filepath"
	path2 "sailing.cn/v2/path"
	"sailing.cn/v2/utils/warning"
	"strings"
)

type Config struct {
	Emqx struct {
		Server   string `json:"server"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"emqx"`
}

var cnf = &Config{}

func (c *Config) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(path2.GetExecPath(), "conf.d", "conf.yml")
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

type Client struct {
	http.Client
	*http.Request
}

// ListResources 获取资源列表
func ListResources() []*Resource {
	c := create("/api/v4/resources", "GET", "").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		log.Errorf("获取资源列表失败:%s", err.Error())
		return nil
	}
	if response.StatusCode == http.StatusOK {
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("获取资源列表失败:%s", err.Error())
			return nil
		}
		list := make([]*Resource, 0)
		result := &Result{}
		json.Unmarshal(buf, result)
		list = result.Data.([]*Resource)
		return list
	}
	return nil
}

// Get 获取单个资源
func (r *Resource) Get() (*Resource, error) {
	c := create("/api/v4/resource/", "GET", "resource:"+r.Id)
	response, err := c.Do(c.Request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, warning.New("获取资源实例失败")
	}
	buf, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

// Create 创建资源
func (r *Resource) Create() error {
	buf, _ := json.Marshal(r)
	c := create("/api/v4/resources", "POST", string(buf)).auth()
	response, err := c.Do(c.Request)
	if err != nil {
		log.Errorf("获取资源列表失败:%s", err.Error())
		return nil
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("创建资源失败:%s", err.Error())
		return err
	}
	buf, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Errorf("获取资源列表失败:%s", err.Error())
		return nil
	}
	return json.Unmarshal(buf, r)
}

// Update 更新资源(api好像没提供直接修改的方法，需删除再创建)
func (r *Resource) Update() error {

	err := r.Delete()
	if err != nil {
		return err
	}
	err = r.Create()
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除资源
func (r *Resource) Delete() error {

	c := create("/api/v4/resources/"+r.Id, "DELETE", "").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return warning.New("删除资源失败")
	}
	buf, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, r)
	return err
}

// Create 创建规则
func (r *Rule) Create() error {
	//form := &url.Values{}
	buf, _ := json.Marshal(r)
	//err := json.Unmarshal(buf, form)
	//if err != nil {
	//	log.Errorf("错误信息：%s", err.Error())
	//}
	//c := create("/api/v4/rules", "POST", *form).auth()
	c := create("/api/v4/rules", "POST", string(buf)).auth()
	response, err := c.Do(c.Request)
	if err != nil {
		log.Errorf("获取规则列表失败：%s", err.Error())
		return err
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("创建规则失败")
		return warning.New(response.Status)
	}
	buf, err = ioutil.ReadAll(response.Body)
	//这里需要关闭
	response.Body.Close()
	if err != nil {
		log.Errorf("获取规则列表失败:%s", err.Error())
		return err
	}
	err = json.Unmarshal(buf, r)
	return err
}

// Update 更新规则
func (r *Rule) Update() error {
	buf, _ := json.Marshal(r)
	c := create("/api/v4/rules", "PUT", string(buf)).auth()
	response, err := c.Do(c.Request)
	if err != nil {
		log.Errorf("修改规则失败")
		return warning.New("修改规则失败")
	}
	if response.StatusCode != http.StatusOK {
		log.Errorf("修改规则失败")
		return warning.New("修改规则失败")
	}
	buf, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return warning.New("修改规则失败")
	}
	err = json.Unmarshal(buf, r)
	return err
}

// UpdateActions 更新规则的响应动作
func (r *Rule) UpdateActions() error {
	buf, err := json.Marshal(r.Actions)
	if err != nil {
		return err
	}
	c := create("/api/v4/rules/"+r.Id, "PUT", "{"+"\"actions\""+":"+string(buf)+"}").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return warning.New("修改规则动作失败")
	}
	buf, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, r)
	return err

}

func (r *Rule) UpdateStatus() error {
	buf, err := json.Marshal(r.Enabled)
	if err != nil {
		return err
	}
	c := create("/api/v4/rules/"+r.Id, "PUT", "{"+"\"enabled\""+":"+string(buf)+"}").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return warning.New("更新规则状态失败")
	}
	buf, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, r)
	return err
}

// Delete 删除规则
func (r *Rule) Delete() error {
	c := create("/api/v4/rules/"+r.Id, "DELETE", "").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return warning.New("删除规则失败")
	}
	buf, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, r)
	return err
}
func create(uri string, method string, values string) *Client {
	c := &Client{}
	if len(cnf.Emqx.Server) <= 0 {
		cnf.Init("")
	}
	uri = cnf.Emqx.Server + uri
	c.Request, _ = http.NewRequest(method, uri, strings.NewReader(values))
	return c
}

// Get 获取单个规则
func (r *Rule) Get() (*Rule, error) {
	c := create("/api/v4/rule/", "GET", "rule:"+r.Id)
	response, err := c.Do(c.Request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, warning.New("获取规则实例失败")
	}
	buf, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

// ListRules 获取资源列表
func ListRules() []*Rule {
	c := create("/api/v4/rule", "GET", "").auth()
	response, err := c.Do(c.Request)
	if err != nil {
		log.Errorf("获取规则列表失败:%s", err.Error())
		return nil
	}
	if response.StatusCode == http.StatusOK {
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Errorf("获取规则列表失败:%s", err.Error())
			return nil
		}
		list := make([]*Rule, 0)
		result := &Result{}
		json.Unmarshal(buf, result)
		list = result.Data.([]*Rule)
		return list
	}
	return nil
}

func (c *Client) auth() *Client {
	if len(cnf.Emqx.Server) <= 0 || len(cnf.Emqx.Password) <= 0 || len(cnf.Emqx.User) < 0 {
		cnf.Init("")
	}
	c.Request.SetBasicAuth(cnf.Emqx.User, cnf.Emqx.Password)
	//c.Request.SetBasicAuth("admin", "public")
	return c
}
