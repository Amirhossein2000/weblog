package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

const (
	AdminRole = iota + 10
	UserRole
)

var secretKey = []byte("SecretKey")

type AuthenticatedUser struct {
	Id   uint
	Role byte
}

type CustomClaims struct {
	jwt.StandardClaims
	*AuthenticatedUser
}

func Authenticate(r *http.Request) *AuthenticatedUser {
	strToken := r.Header.Get("Authorization")
	if strToken == "" {
		return nil
	}

	claims := CustomClaims{}
	token, err := jwt.ParseWithClaims(strToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Println("JWT err:", err.Error())
		return nil
	}

	if !token.Valid {
		return nil
	}

	if claims.AuthenticatedUser != nil {
		return claims.AuthenticatedUser
	}

	return nil
}

func HasPermission(userID uint, authUser *AuthenticatedUser) bool {
	return userID == authUser.Id || authUser.Role == AdminRole
}

func GenerateToken(user *AuthenticatedUser) string {
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Subject:   "authentication",
		},

		AuthenticatedUser: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err := token.SignedString(secretKey)

	if err != nil {
		panic(err)
	}

	return strToken
}
