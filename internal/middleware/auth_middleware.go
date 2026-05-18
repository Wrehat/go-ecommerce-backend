package middleware

import (
	"ecommerce/pkg/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(jwtsecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Todo : cek request bawa token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.JSON(c, http.StatusUnauthorized, "Akses ditolak. Token tidak ditemukan", nil)
			c.Abort()
			return
		}

		// Todo : cek format token valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.JSON(c, http.StatusUnauthorized, "Format token tidak valid", nil)
			c.Abort()
			return
		}

		// Todo : cek keaslian token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			// Todo : cek signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritma token tidak valid")
			}

			return []byte(jwtsecret), nil
		})

		// Todo : handle error
		// Jika stempel palsu, rusak, atau token expired
		if err != nil || !token.Valid {
			response.JSON(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa", nil)
			c.Abort()
			return
		}

		// Todo : ekstrak token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
		} else {
			response.JSON(c, http.StatusUnauthorized, "Gagal membaca token", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Todo : ambil role
		userRole, exist := c.Get("role")
		if !exist {
			response.JSON(c, http.StatusForbidden, "Akses ditolak. Role tidak dikenali", nil)
			c.Abort()
			return
		}

		// Todo : validasi role
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		// Todo : atur hak akses
		if !isAllowed {
			response.JSON(c, http.StatusForbidden, "Akses ditolak. Anda tidak memiliki hak akses untuk halaman ini", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
