package jwtutil

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	UserID               int    `json:"user_id"`
	Email                string `json:"email"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Wajib dimasukkan untuk menampung data bawaan JWT (seperti exp, iat, dll)
}
