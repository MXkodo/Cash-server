package middleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	Secret []byte
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{[]byte(secret)}
}

func (m *JWTMiddleware) UserIDFromTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization") 

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return m.Secret, nil
		})

		if err != nil || !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
			c.Abort()
			return
		}

		if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("userID", int(userID)) 
			}

			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен истек"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
