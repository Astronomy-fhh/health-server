package controller_web

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO 环境变量
var secretKey = "7bee5d7382c34f7af125efe692417016ed4456fdc48392wef7b865364317b1bc419b477c"

type TokenPayload struct {
	AdminId uint64 `json:"admin_id"`
}

const Exp = time.Hour * 6

func GenerateAdminToken(payload TokenPayload, exp int64) (string, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"payload": payloadJson, // 将整个 payload 存储在 claims 中
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseAdminToken(tokenString string) (string, bool, error) {
	expired := false
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				expired = true
			}
		}

		payloadJSON := claims["payload"].(string)
		return payloadJSON, expired, nil
	}

	return "", false, fmt.Errorf("invalid token")
}
