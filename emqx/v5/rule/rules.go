package rule

import (
	"errors"
	"fmt"
	"net/http"
	"sailing.cn/emqx/v5/client"
	"sailing.cn/utils"
)

const (
	rulesURL = "/rules"
)

type RuleService client.Service

type EMQXRule struct {
	// The ID of the rule
	Id      string      `json:"id"`
	Metrics interface{} `json:"metrics,omitempty"`
	// The metrics of the rule for each node
	NodeMetrics []interface{} `json:"node_metrics,omitempty"`
	// The topics of the rule
	From []string `json:"from,omitempty"`
	// The created time of the rule
	CreatedAt string `json:"created_at,omitempty"`
	// The name of the rule
	Name string `json:"name,omitempty"`
	// </br>SQL query to transform the messages.</br>Example: <code>SELECT * FROM \"test/topic\" WHERE payload.x = 1</code></br>
	Sql string `json:"sql"`
	// </br>A list of actions of the rule.</br>An action can be a string that refers to the channel ID of an EMQX bridge, or an object</br>that refers to a function.</br>There a some built-in functions like \"republish\" and \"console\", and we also support user</br>provided functions in the format: \"{module}:{function}\".</br>The actions in the list are executed sequentially.</br>This means that if one of the action is executing slowly, all the following actions will not</br>be executed until it returns.</br>If one of the action crashed, all other actions come after it will still be executed, in the</br>original order.</br>If there's any error when running an action, there will be an error message, and the 'failure'</br>counter of the function action or the bridge channel will increase.</br>
	Actions []interface{} `json:"actions,omitempty"`
	// Enable or disable the rule
	Enable bool `json:"enable,omitempty"`
	// The description of the rule
	Description string `json:"description,omitempty"`
	// Rule metadata, do not change manually
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EMQXRuleCreateRequest 规则创建模型
type EMQXRuleCreateRequest struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name,omitempty"`
	Sql         string                   `json:"sql"`
	Actions     []map[string]interface{} `json:"actions,omitempty"`
	Enable      bool                     `json:"enable,omitempty"`
	Description string                   `json:"description,omitempty"`
	Metadata    map[string]interface{}   `json:"metadata,omitempty"`
}

// EMQXRuleUpdateRequest 规则修改模型
type EMQXRuleUpdateRequest struct {
	Name string `json:"name,omitempty"`
	Sql  string `json:"sql"`
	//Actions []map[string]interface{} `json:"actions,omitempty"`
	Actions     []interface{}          `json:"actions"` //这里不需要省略参数，会出现删除最后一个动作失败
	Enable      bool                   `json:"enable,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type RuleQuery struct {
}

// GetRule 查询规则
func (c *RuleService) GetRule(ruleId string) (rule *EMQXRule, err error) {
	result := &EMQXRule{}
	r := c.Client.R().
		SetResult(result)
	resp, err := r.Get(c.Client.RequestURL(rulesURL + "/" + ruleId))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return rule, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	rule = result
	return
}

// ListRule 规则列表
func (c *RuleService) ListRule(query *RuleQuery) (*client.PageResponse[EMQXRule], error) {
	result := new(client.PageResponse[EMQXRule])
	_query := utils.ToMapStr(*query)
	resp, err := c.Client.R().
		SetQueryParams(_query).
		SetResult(&result).Get(c.Client.RequestURL(rulesURL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	return result, nil
}

// CreateRule 创建规则
func (c *RuleService) CreateRule(rule *EMQXRuleCreateRequest) (interface{}, error) {
	var result interface{}
	resp, err := c.Client.R().SetBody(rule).SetResult(&result).Post(c.Client.RequestURL(rulesURL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("请求错误:code%d msg::%s", resp.StatusCode, resp.String()))
	}
	return result, nil
}

// UpdateEnable 启用/停用 规则
func (c *RuleService) UpdateEnable(ruleId string, enable bool) (interface{}, error) {
	var result interface{}
	url := c.Client.RequestURL(rulesURL + "/" + ruleId)
	resp, err := c.Client.R().SetBody(struct {
		Enable bool `json:"enable"`
	}{enable}).SetResult(&result).Put(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("请求错误:code%d msg::%s", resp.StatusCode, resp.String()))
	}
	return result, nil
}

// DeleteRule 删除规则
func (c *RuleService) DeleteRule(ruleId string) error {
	resp, err := c.Client.R().Delete(c.Client.RequestURL(rulesURL + "/" + ruleId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("EMQX删除规则出现错误:" + resp.String())
	}
	return nil
}

// Update 更新转发规则
func (c *RuleService) Update(ruleId string, rule *EMQXRuleUpdateRequest) (interface{}, error) {
	model, err := c.GetRule(ruleId)
	if err != nil {
		return nil, errors.New("当前规则未找到，请重试")
	}
	//赋值部分数据
	rule.Name = model.Name
	rule.Sql = model.Sql
	rule.Enable = model.Enable
	rule.Description = model.Description
	rule.Metadata = model.Metadata
	//xxx, err := json.Marshal(rule)
	//log.Error(string(xxx))
	var result interface{}
	resp, err := c.Client.R().SetBody(rule).SetResult(&result).Put(c.Client.RequestURL(rulesURL + "/" + ruleId))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("请求错误:code%d msg::%s", resp.StatusCode, resp.String()))
	}
	return result, nil
}
