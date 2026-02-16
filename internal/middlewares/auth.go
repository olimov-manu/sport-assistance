package middlewares

import (
	"fmt"
	"net/http"
	"sport-assistance/pkg/myerrors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      uint64   `json:"user_id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Error(myerrors.AuthorizationHeaderEmptyErrorMessage)
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: myerrors.AuthorizationHeaderEmptyErrorMessage})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			m.logger.Info(myerrors.InvalidBearerTokenFormatErrorMessage)
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: myerrors.InvalidBearerTokenFormatErrorMessage})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
			}

			return []byte(m.cfg.AccessTokenSecret), nil // если AccessTokenSecret строка
		})

		if err != nil {
			m.logger.Error("jwt parse error: %v", err)
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: myerrors.ParseTokenErrorMessage})
			c.Abort()
			return
		}

		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: myerrors.InvalidTokenErrorMessage})
			c.Abort()
			return
		}

		if claims.Email == "" {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: "Token does not contain user email"})
			c.Abort()
			return
		}

		isUserExists, err := m.repo.UserExistsByEmail(ctx, claims.Email)
		if isUserExists == false {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: "Пользователя с таким email не существует."})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: myerrors.CheckUserExistsByEmailErrorMessage})
			c.Abort()
			return
		}

		key := fmt.Sprintf(m.cfg.AccessTokenRedisPrefix, claims.UserID)

		ttl, err := m.redisClient.TTL(ctx, key).Result()
		if err != nil || ttl.Seconds() < 1 {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: "Invalid or expired token"})
			c.Abort()
			return
		}

		m.logger.Info(myerrors.CheckUserExistsByEmailErrorMessage)

		if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
			c.JSON(http.StatusUnauthorized, AuthResponse{Success: false, Error: "Token expired"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("permissions", claims.Permissions)
		c.Set("claims", claims)
		c.Next()
	}
}
