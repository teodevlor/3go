package controller

import (
	"go-structure/internal/common"
	"go-structure/internal/middleware"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	BaseController struct{}

	ClientInfo struct {
		IP        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		RequestID string `json:"request_id"`
	}
)

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (b *BaseController) GetAccountIDFromContext(c *gin.Context) (uuid.UUID, *common.ResponseData) {
	accountIDValue, exists := c.Get(middleware.ContextAccountIDKey)
	if !exists {
		return uuid.Nil, common.ErrorResponse(common.StatusUnauthorized, []string{"unauthorized"})
	}

	accountID, ok := accountIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, common.ErrorResponse(common.StatusUnauthorized, []string{"invalid account id in context"})
	}

	return accountID, nil
}

func (b *BaseController) GetClientIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

func (b *BaseController) GetUserAgent(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		return "Unknown"
	}
	return userAgent
}

func (b *BaseController) GetRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

// GetBearerToken lấy Bearer token từ Authorization header
func (b *BaseController) GetBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func (b *BaseController) GetClientInfo(c *gin.Context) ClientInfo {
	return ClientInfo{
		IP:        b.GetClientIP(c),
		UserAgent: b.GetUserAgent(c),
		RequestID: b.GetRequestID(c),
	}
}
