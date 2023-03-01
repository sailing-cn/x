package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"sailing.cn/zentao/cache"
	"sync"
	"time"
)

const (
	tokenURL = "/tokens"
)

type Token struct {
	User      string
	Password  string
	domain    string
	cache     cache.Cache
	tokenLock *sync.Mutex
}
type AccessToken struct {
	Token string `json:"token"`
}

type BasicAuth struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
type TokenHandle interface {
	GetToken() (token string, err error)
	GetTokenContext(ctx context.Context) (token string, err error)
}

func NewToken(user, password string, domain string, cache cache.Cache) TokenHandle {
	return &Token{
		User:      user,
		Password:  password,
		domain:    domain,
		tokenLock: new(sync.Mutex),
		cache:     cache,
	}
}

func (t *Token) GetToken() (token string, err error) {
	return t.GetTokenContext(context.Background())
}

func (t *Token) GetTokenContext(ctx context.Context) (token string, err error) {
	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()
	cacheKey := fmt.Sprintf("s_access_token_%s", t.User)
	val := t.cache.Get(cacheKey)
	if val != nil {
		token = val.(string)
		return
	}
	var resToken AccessToken
	body := BasicAuth{Account: t.User, Password: t.Password}
	resToken, err = getTokenFormServerContext(t.domain+tokenURL, body)
	if err != nil {
		return
	}
	expires := time.Duration(1440) * time.Second
	err = t.cache.Set(cacheKey, resToken, expires)
	token = resToken.Token
	return
}

func getTokenFormServerContext(url string, data interface{}) (token AccessToken, err error) {
	//var body []byte
	//jsonBuf := new(bytes.Buffer)
	//enc := json.NewEncoder(jsonBuf)
	//enc.SetEscapeHTML(false)
	//err = enc.Encode(data)
	//if err != nil {
	//	return
	//}
	//response, err := http.Post(url, "application/json;charset=utf-8", jsonBuf)
	//if err != nil {
	//	return
	//}
	//defer response.Body.Close()
	//
	//if response.StatusCode != http.StatusCreated {
	//	err = errors.New(fmt.Sprintf("http get error : uri=%v , statusCode=%v", url, response.StatusCode))
	//	return
	//}
	//body, err = io.ReadAll(response.Body)
	//if err != nil {
	//	return
	//}
	//err = json.Unmarshal(body, &token)
	//if err != nil {
	//	return
	//}
	token = AccessToken{}
	client := req.C()
	response, err := client.R().SetBody(data).SetResult(&token).Post(url)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusCreated {
		err = errors.New(fmt.Sprintf("http get error : uri=%v , statusCode=%v", url, response.StatusCode))
		return
	}
	return
}
