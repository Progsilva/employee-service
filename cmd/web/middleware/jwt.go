package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"userID"`
	jwt.RegisteredClaims
}

var UserKey = "user-key"

type User struct {
	ID int64 `json:"id"`
}

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateToken(c, secret); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func validateToken(c *gin.Context, secret string) error {
	token, err := getToken(c, secret)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), UserKey, &User{ID: claims.UserID}))
		return nil
	}
	return errors.New("invalid token provided")
}

func getToken(c *gin.Context, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(getTokenFromRequest(c), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func CtxUser(ctx context.Context) *User {
	raw, _ := ctx.Value(UserKey).(*User)
	return raw
}
