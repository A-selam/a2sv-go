package infrastructure

import (
	"fmt"
	domain "task_manager/Domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jWTService struct {
	secret []byte
}

func NewJWTService(secret string) domain.JWTService {
	return &jWTService{secret: []byte(secret)}
}

func (j *jWTService) GenerateToken(username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(j.secret))
}


func (s *jWTService) ParseToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	username, ok1 := claims["username"].(string)
	role, ok2 := claims["role"].(string)

	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("invalid token payload")
	}

	return username, role, nil
}


