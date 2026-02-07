package middleware

import (
	"strings"

	"go-structure/internal/common"
	jwtutil "go-structure/internal/utils/jwt"

	"github.com/gin-gonic/gin"
)

const ContextAccountIDKey = "account_id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			resp := common.ErrorResponse(common.StatusUnauthorized, []string{"Thiếu hoặc sai header Authorization"})
			c.AbortWithStatusJSON(common.StatusUnauthorized, resp)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		accountID, err := jwtutil.ParseAccessToken(tokenStr)
		if err != nil {
			resp := common.ErrorResponse(common.StatusUnauthorized, []string{"Token không hợp lệ hoặc đã hết hạn"})
			c.AbortWithStatusJSON(common.StatusUnauthorized, resp)
			return
		}

		c.Set(ContextAccountIDKey, accountID)
		c.Next()
	}
}
