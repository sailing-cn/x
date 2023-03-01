package project

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"zentao/model"
)

const (
	projectURL = "projects"
)

type ProjectPageResult struct {
	model.PageResult
	Projects []Project `json:"projects"`
}

type Project struct {
	ID            int         `json:"id"`
	Project       int         `json:"project"`
	Model         string      `json:"model"`
	Type          string      `json:"type"`
	Lifetime      string      `json:"lifetime"`
	Budget        string      `json:"budget"`
	BudgetUnit    string      `json:"budgetUnit"`
	Attribute     string      `json:"attribute"`
	Percent       int         `json:"percent"`
	Milestone     string      `json:"milestone"`
	Output        string      `json:"output"`
	Auth          string      `json:"auth"`
	Parent        int         `json:"parent"`
	Path          string      `json:"path"`
	Grade         int         `json:"grade"`
	Name          string      `json:"name"`
	Code          string      `json:"code"`
	Begin         string      `json:"begin"`
	End           string      `json:"end"`
	RealBegan     string      `json:"realBegan"`
	RealEnd       interface{} `json:"realEnd"`
	Days          int         `json:"days"`
	Status        string      `json:"status"`
	SubStatus     string      `json:"subStatus"`
	Pri           string      `json:"pri"`
	Desc          string      `json:"desc"`
	Version       int         `json:"version"`
	ParentVersion int         `json:"parentVersion"`
	PlanDuration  int         `json:"planDuration"`
	RealDuration  int         `json:"realDuration"`
	OpenedBy      struct {
		ID       int    `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"openedBy"`
	OpenedDate    time.Time `json:"openedDate"`
	OpenedVersion string    `json:"openedVersion"`
	LastEditedBy  struct {
		ID       int    `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"lastEditedBy"`
	LastEditedDate time.Time   `json:"lastEditedDate"`
	ClosedBy       interface{} `json:"closedBy"`
	ClosedDate     interface{} `json:"closedDate"`
	CanceledBy     interface{} `json:"canceledBy"`
	CanceledDate   interface{} `json:"canceledDate"`
	SuspendedDate  string      `json:"suspendedDate"`
	PO             string      `json:"PO"`
	PM             struct {
		ID       int    `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"PM"`
	QD           string        `json:"QD"`
	RD           string        `json:"RD"`
	Team         string        `json:"team"`
	ACL          string        `json:"acl"`
	Whitelist    []interface{} `json:"whitelist"`
	Order        int           `json:"order"`
	Vision       string        `json:"vision"`
	DisplayCards int           `json:"displayCards"`
	FluidBoard   string        `json:"fluidBoard"`
	Deleted      bool          `json:"deleted"`
	Hours        struct {
		TotalEstimate int `json:"totalEstimate"`
		TotalConsumed int `json:"totalConsumed"`
		TotalLeft     int `json:"totalLeft"`
		Progress      int `json:"progress"`
		TotalReal     int `json:"totalReal"`
	} `json:"hours"`
	TeamCount     int           `json:"teamCount"`
	LeftTasks     string        `json:"leftTasks"`
	TeamMembers   []interface{} `json:"teamMembers"`
	TotalEstimate int           `json:"totalEstimate"`
	TotalConsumed int           `json:"totalConsumed"`
	TotalLeft     int           `json:"totalLeft"`
	Progress      int           `json:"progress"`
	TotalReal     int           `json:"totalReal"`
}

func (c *Client) ProjectList(page, limit int) (interface{}, error) {
	result := &ProjectPageResult{}
	resp, err := c.ctx.Client.R().
		SetHeader("Token", c.ctx.Token).
		SetQueryParams(map[string]string{
			"page":  fmt.Sprintf("%d", page),
			"limit": fmt.Sprintf("%d", limit),
		}).SetResult(result).Get(c.ctx.RequestURL(projectURL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	return result, nil
}
