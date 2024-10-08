package kit

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO 环境变量
var secretKey = "7bee5d738c34f7af5efe6917016ed4456fdc4839f7b865364317b1bc419b477c"

type UserTokenPayload struct {
	Uid string `json:"uid"`
}

func GenerateUserToken(payload *UserTokenPayload, exp int64) (string, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"payload": string(payloadJSON), // 将整个 payload 存储在 claims 中
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseUserToken(tokenString string) (*UserTokenPayload, bool, error) {
	expired := false
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				expired = true
			}
		}

		payloadJSON := claims["payload"].(string)

		var payload UserTokenPayload
		err := json.Unmarshal([]byte(payloadJSON), &payload)
		if err != nil {
			return nil, false, err
		}
		return &payload, expired, nil
	}

	return nil, false, fmt.Errorf("invalid token")
}
