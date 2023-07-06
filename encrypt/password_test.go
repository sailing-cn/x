package encrypt

import (
	"testing"
)

func TestSalt(t *testing.T) {
	salt, err := Salt()
	if err != nil {
		t.Fail()
		return
	}
	t.Logf("salt:%s", salt)
}

func TestPassword(t *testing.T) {
	pwd := "123456"
	//salt, err := Salt()
	//if err != nil {
	//	t.Errorf("生成盐出错:%s", err.Error())
	//	return
	//}
	password, err := encryptPassword(pwd)
	if err != nil {
		t.Errorf("密码加密出错:%s", err.Error())
		return
	}
	//t.Logf("盐:%s", hex.EncodeToString(salt))
	t.Logf("密:%s", string(password))
}
