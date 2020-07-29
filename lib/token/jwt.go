package token

import (
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type CustomClaims struct {
	UserName string `json:"user_name"` //这里放自定义的struct或字段
	jwt.StandardClaims
}

// 这个是jwt默认的几个字段
// type StandardClaims struct {
// 	Audience  string `json:"aud,omitempty"`
// 	ExpiresAt int64  `json:"exp,omitempty"`
// 	Id        string `json:"jti,omitempty"`
// 	IssuedAt  int64  `json:"iat,omitempty"`
// 	Issuer    string `json:"iss,omitempty"`
// 	NotBefore int64  `json:"nbf,omitempty"`
// 	Subject   string `json:"sub,omitempty"`
// }

// Token jwt服务
type Token struct {
	rwlock     sync.RWMutex
	privateKey []byte
}

func (srv *Token) get() []byte {
	srv.rwlock.RLock()
	defer srv.rwlock.RUnlock()

	return srv.privateKey
}

func (srv *Token) put(newKey []byte) {
	srv.rwlock.Lock()
	defer srv.rwlock.Unlock()

	srv.privateKey = newKey
}
func (srv *Token) Init(newKey []byte) {
	srv.rwlock.Lock()
	defer srv.rwlock.Unlock()

	if len(newKey) > 0 {
		srv.privateKey = newKey
	}
}

//Decode 解码
func (srv *Token) Decode(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.get(), nil
	})

	if err != nil {
		return nil, err
	}
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}

// Encode 将 User 用户信息加密为 JWT 字符串
// expireTime := time.Now().Add(time.Hour * 24 * 3).Unix() 三天后过期
func (srv *Token) Encode(issuer, userName string, expireTime int64) (string, error) {
	claims := CustomClaims{
		userName, //这里可以替换放自定义的struct
		jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(srv.get())
}
