package jwt

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

var JWKS *keyfunc.JWKS

func InitJWKS(url string) {
	if len(url) <= 0 {
		log.Fatalf("认证地址未配置")
		return
	}
	url = fmt.Sprintf("%s/.well-known/openid-configuration/jwks", url)
	ctx, cancel := context.WithCancel(context.Background())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
		Client:            client,
	}
	jwks, err := keyfunc.Get(url, options)
	if err != nil {
		log.Fatalf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
	}
	cancel()
	jwks.EndBackground()
	JWKS = jwks
}

func Auth(c *gin.Context) {
	// 从HTTP Header中获取Token
	tokenString := c.GetHeader("Authorization")
	const bearerLength = len("Bearer ")
	if len(tokenString) < bearerLength {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "header Authorization has not Bearer token"})
		return
	}
	tokenString = strings.TrimSpace(tokenString[bearerLength:])

	// 解析Token
	token, err := jwt.Parse(tokenString, JWKS.Keyfunc)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// 验证Token是否有效
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	// Token验证通过，将Token中的用户信息存储到Gin的Context中
	claims := token.Claims.(jwt.MapClaims)
	for _, claim := range CLAIMS {
		c.Set(claim, claims[claim])
	}
}
