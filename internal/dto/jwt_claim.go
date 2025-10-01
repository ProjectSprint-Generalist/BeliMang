package dto

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthUser struct {
	Username string
	Email    string
	Role     string
}
