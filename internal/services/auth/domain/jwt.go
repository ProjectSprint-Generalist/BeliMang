package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	UserID   string
	Username string
	Email    string
	Role     string
	jwt.RegisteredClaims
}
