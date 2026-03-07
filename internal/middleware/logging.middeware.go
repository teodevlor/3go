package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"go-structure/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	KeyCorrelationIDHeaderKey = "X-Correlation-ID"
	ContextCorrelationIDKey   = "correlation_id"
	maskedValue               = "***"
)

var sensitiveKeys = map[string]struct{}{
	"password":         {},
	"new_password":     {},
	"old_password":     {},
	"confirm_password": {},
	"token":            {},
	"access_token":     {},
	"refresh_token":    {},
	"secret":           {},
	"secret_key":       {},
	"credit_card":      {},
	"card_number":      {},
	"cvv":              {},
	"pin":              {},
	"otp":              {},
}

type correlationIDContextKey struct{}

var correlationIDKey = correlationIDContextKey{}

func CorrelationIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(correlationIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// maskSensitiveFields đệ quy duyệt map JSON và mask các key nhạy cảm.
func maskSensitiveFields(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))
	for k, v := range data {
		if _, isSensitive := sensitiveKeys[strings.ToLower(k)]; isSensitive {
			result[k] = maskedValue
			continue
		}
		switch val := v.(type) {
		case map[string]any:
			result[k] = maskSensitiveFields(val)
		case []any:
			result[k] = maskSlice(val)
		default:
			result[k] = v
		}
	}
	return result
}

func maskSlice(arr []any) []any {
	result := make([]any, len(arr))
	for i, item := range arr {
		if m, ok := item.(map[string]any); ok {
			result[i] = maskSensitiveFields(m)
		} else {
			result[i] = item
		}
	}
	return result
}

func parseAndMaskJSON(body []byte) map[string]any {
	if len(body) == 0 {
		return nil
	}
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil
	}
	return maskSensitiveFields(parsed)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		correlationID := c.GetHeader(KeyCorrelationIDHeaderKey)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		// 2. Gắn vào gin.Context và response header
		c.Set(ContextCorrelationIDKey, correlationID)
		c.Writer.Header().Set(KeyCorrelationIDHeaderKey, correlationID)

		// 3. Gắn vào context chuẩn của request
		ctx := context.WithValue(c.Request.Context(), correlationIDKey, correlationID)
		c.Request = c.Request.WithContext(ctx)

		// 4. Đọc request body (đọc xong phải gắn lại để handler vẫn dùng được)
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		rawQuery := c.Request.URL.RawQuery
		endpoint := c.FullPath()

		// 5. Xử lý request
		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 6. Log tất cả request: status >= 400 thì kèm body để debug, còn lại chỉ log lifecycle
		if status >= 400 {
			maskedReqBody := parseAndMaskJSON(reqBodyBytes)
			reqBodyField := zap.Skip()
			if maskedReqBody != nil {
				reqBodyField = zap.Any("request_body", maskedReqBody)
			}
			global.Logger.Warn("http_request",
				zap.String("correlation_id", correlationID),
				zap.String("method", method),
				zap.String("endpoint", endpoint),
				zap.String("query", rawQuery),
				zap.Int("status", status),
				zap.String("client_ip", clientIP),
				zap.String("duration_ms", fmt.Sprintf("%dms", latency.Milliseconds())),
				zap.String("user_agent", userAgent),
				reqBodyField,
			)
		} else {
			global.Logger.Info("http_request",
				zap.String("correlation_id", correlationID),
				zap.String("method", method),
				zap.String("endpoint", endpoint),
				zap.String("query", rawQuery),
				zap.Int("status", status),
				zap.String("client_ip", clientIP),
				zap.String("duration_ms", fmt.Sprintf("%dms", latency.Milliseconds())),
				zap.String("user_agent", userAgent),
			)
		}
	}
}
