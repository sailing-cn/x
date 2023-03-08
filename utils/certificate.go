package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
)

// GetCertificateInfo 获取证书的基本信息
func GetCertificateInfo(content string) CertificateInfo {
	result := CertificateInfo{}
	x509Cert, err := x509.ParseCertificate(IsLegalBytes([]byte(content)))
	if err != nil {
		log.Error("解析证书文件失败")
		return result
	}
	result.CommonName = x509Cert.Issuer.CommonName
	result.VerifyCode = x509Cert.Subject.CommonName
	result.Owner = x509Cert.Issuer.String()
	result.EffectiveDate = x509Cert.NotBefore.UnixMilli()
	result.ExpiryDate = x509Cert.NotAfter.UnixMilli()
	result.Md5 = GetMd5Value(content)
	return result
}

// VerityCertificate  验证证书
func VerityCertificate(certPEM string, rootPEM []string) (result bool, sha1 string, sha256 string) {
	result = false
	bytes, err := base64.StdEncoding.DecodeString(certPEM)
	if err != nil {
		return false, "", ""
	}
	roots := x509.NewCertPool()
	cert, err := x509.ParseCertificate(bytes)
	if err != nil {
		log.Errorf("failed to parse certificate : " + err.Error())
		return result, sha1, sha256
	}
	for _, s := range rootPEM {
		ok := roots.AppendCertsFromPEM([]byte(s))
		if !ok {
			log.Error("failed to parse root certificate")
			return result, sha1, sha256
		}
	}
	opts := x509.VerifyOptions{
		Roots: roots,
	}
	if _, err := cert.Verify(opts); err != nil {
		log.Errorf("failed to verify certificate: " + err.Error())
		return result, sha1, sha256
	}

	sha1, sha256 = GetFingerprint(bytes)

	return true, sha1, sha256
}

// VerityCertificateIncludeCode  验证证书
func VerityCertificateIncludeCode(certPEM string, rootPEM string, verityCode string) bool {
	roots := x509.NewCertPool()
	cert, err := x509.ParseCertificate(IsLegalBytes([]byte(certPEM)))
	if err != nil {
		log.Errorf("failed to parse certificate : " + err.Error())
		return false
	}
	// 验证验证码是否正确
	if len(verityCode) <= 0 || cert.Subject.CommonName != verityCode {
		return false
	}
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Error("failed to parse root certificate")
		return false
	}
	opts := x509.VerifyOptions{
		Roots: roots,
	}
	if _, err := cert.Verify(opts); err != nil {
		log.Errorf("failed to verify certificate: " + err.Error())
		return false
	}
	return true
}

type CertificateInfo struct {
	CommonName    string `json:"common_name"`    // 公共名字
	EffectiveDate int64  `json:"effective_date"` // 证书激活时间
	ExpiryDate    int64  `json:"expiry_date"`    // 证书失效原因
	Owner         string `json:"owner"`          // 证书拥有者
	VerifyCode    string `json:"verify_code"`    // 验证码
	Md5           string `json:"md_5"`           //证书标识
}

// GetFingerprint 获取证书的 sha_1 和 sha_256 指纹 (emqx  base 解码得到的)
func GetFingerprint(content []byte) (sha_1 string, sha_256 string) {
	certificate, err := x509.ParseCertificate(content)
	if err != nil {
		log.Errorf("failed to parse certificate : " + err.Error())
		return sha_1, sha_256
	}
	h := sha1.New()
	h.Write(certificate.Raw)
	sha_1 = hex.EncodeToString(h.Sum(nil))
	h1 := sha256.New()
	h1.Write(certificate.Raw)
	sha_256 = hex.EncodeToString(h1.Sum(nil))
	return
}

// IsLegal 是否为证书格式
func IsLegal(data []byte) bool {
	block, _ := pem.Decode(data)
	if block == nil {
		return false
	}
	return true
}

// IsLegalBytes 返回证书bytes
func IsLegalBytes(data []byte) []byte {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil
	}
	return block.Bytes
}

func GetMd5Value(content string) string {
	certificate, err := x509.ParseCertificate(IsLegalBytes([]byte(content)))
	if err != nil {
		return ""
	}
	m := md5.New()
	m.Write(certificate.Raw)
	return hex.EncodeToString(m.Sum(nil))
}
