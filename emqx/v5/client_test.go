package v5

import (
	"sailing.cn/v2/emqx/v5/client"
	"sailing.cn/v2/emqx/v5/config"
	"testing"
	"time"
)

var (
	local = &config.Configuration{
		Server:     "http://192.168.1.65:18083",
		ApiSecret:  "l3vIoWtA4t7lsSyR0bkA9AfVR9BZ1T228RCsZ2Q2xmTgF",
		ApiKey:     "d49ed2466159d244",
		ApiVersion: "v5",
	}
	tencent = &config.Configuration{
		Server:     "http://101.33.245.172:28083",
		ApiSecret:  "oEuxKj8zR2JAXfdI9AHKqDsxVyGva029CggIMSm9BOlWZJ",
		ApiKey:     "ac84a121a5677f0f",
		ApiVersion: "v5",
	}
)

func TestListClient(t *testing.T) {
	k1 := time.Now().UnixMilli()
	println(k1)
	cfg := local
	emqx := NewEmqxClientWithCnf(cfg)
	k := client.ClientQuery{LikeClientId: "IOT"}
	list, err := emqx.ClientService.ListClient(&client.ClientPageQuery{
		ClientQuery: k,
	})
	if err != nil {
		t.Errorf(err.Error())
	} else {
		t.Logf("%v", list)
	}
}

func TestKickClient(t *testing.T) {
	cfg := local
	emqx := NewEmqxClientWithCnf(cfg)
	err := emqx.ClientService.KickClient("IOT_PLATFORM_0_0_2023021217")
	if err != nil {
		t.Errorf(err.Error())
	}
}
