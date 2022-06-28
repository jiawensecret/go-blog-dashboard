package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"platform/global"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	ID         uint32
	Username   string
	Mobile     string
	BufferTime int64
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("token已过期")
	TokenNotValidYet = errors.New("token无效")
	TokenMalformed   = errors.New("token无效")
	TokenInvalid     = errors.New("token无效:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.Config.Jwt.SigningKey),
	}
}

func GetRedisJWT(mobile string) (err error, redisJWT string) {
	redisJWT, err = global.Redis.Get("login:" + mobile).Result()
	return err, redisJWT
}

func SetRedisJWT(jwt string, mobile string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.Config.Jwt.ExpiresTime) * time.Second
	err = global.Redis.Set("login:"+mobile, jwt, timer).Err()
	return err
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
