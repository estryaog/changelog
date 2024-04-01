package types

import "github.com/golang-jwt/jwt"

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type Changelog struct {
	Id        string `json:"id"`
	Version   string `json:"version"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type JWTClaims struct {
	Id      string `json:"id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}
