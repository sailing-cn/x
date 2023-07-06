package encrypt

import (
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm3"
)

func Salt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}

func EncryptPassword(salt []byte, password string) ([]byte, error) {
	hash := sm3.New()
	hash.Write([]byte(password))
	pwdHash := hash.Sum(nil)
	hash.Write(pwdHash)
	hash.Write(salt)
	combinedHash := hash.Sum(nil)
	finalHash := make([]byte, len(salt)+len(combinedHash))
	copy(finalHash[:32], salt)
	copy(finalHash[32:], combinedHash)
	return []byte(hex.EncodeToString(finalHash)), nil
}

// 加密密码
func encryptPassword(password string) ([]byte, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	hash := sm3.New()
	hash.Write([]byte(password))
	passwordHash := hash.Sum(nil)

	hash.Write(passwordHash)
	hash.Write(salt)
	combinedHash := hash.Sum(nil)

	finalHash := make([]byte, len(salt)+len(combinedHash))
	copy(finalHash[:8], salt)
	copy(finalHash[8:], combinedHash)
	log.Infof("盐:%s", hex.EncodeToString(salt))
	log.Infof("密:%s", hex.EncodeToString(finalHash))
	return []byte(hex.EncodeToString(finalHash)), nil
}
