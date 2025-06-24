package utils

import (
	db "TodoApp/database"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))
var refreshToken = []byte(os.Getenv("REFRESH_TOKEN"))

func GenerateJWT(userID int) (string, error) {
	//creating clam
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Second * 5).Unix(), //24 hours
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	fmt.Println(tokenStr)
	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired token...")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid sub clam")
	}
	return int(sub), nil
}

func GenerateRefreshToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"typ": "refresh",                                 // type for refresh oken
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days lifespan
	})
	return token.SignedString(refreshToken)
}

func ParseRefreshToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return refreshToken, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claim")
	}
	if claims["typ"] != "refresh" {
		return 0, errors.New("token is not a refresh token")
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid sub claim")
	}
	return int(sub), nil // returinig user_id and null as error
}

func IsRefreshTokenValid(userID int, token string) (bool, error) {
	isExist := false
	err := db.DB.Get(&isExist, `
		SELECT EXISTS (
			SELECT 1 FROM sessions
			WHERE user_id = $1
			AND refresh_token= $2
			AND expires_at > NOW()
			AND archived_at IS NULL
		)
	`, userID, token)

	if err != nil {
		return false, err
	}
	return isExist, nil //return true or 1 and err ans null
}
