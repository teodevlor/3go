package middleware

import (
	"context"
	"go-structure/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminPermissionChecker interface {
	AdminHasPermission(ctx context.Context, adminID uuid.UUID, permissionCode string) (bool, error)
}

func RequirePermission(checker AdminPermissionChecker, permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if checker == nil || len(permissionCodes) == 0 {
			c.Next()
			return
		}

		adminIDVal, exists := c.Get(ContextAdminIDKey)
		if !exists {
			resp := common.ErrorResponse(common.StatusUnauthorized, []string{"Chưa đăng nhập (admin)"})
			c.AbortWithStatusJSON(common.StatusUnauthorized, resp)
			return
		}

		adminID, ok := adminIDVal.(uuid.UUID)
		if !ok {
			resp := common.ErrorResponse(common.StatusUnauthorized, []string{"Định dạng admin không hợp lệ"})
			c.AbortWithStatusJSON(common.StatusUnauthorized, resp)
			return
		}

		for _, code := range permissionCodes {
			has, err := checker.AdminHasPermission(c.Request.Context(), adminID, code)
			if err != nil {
				resp := common.ErrorResponse(common.StatusInternalServerError, []string{"Lỗi kiểm tra quyền"})
				c.AbortWithStatusJSON(common.StatusInternalServerError, resp)
				return
			}
			if has {
				c.Next()
				return
			}
		}

		resp := common.ErrorResponse(common.StatusForbidden, []string{"Bạn không có quyền thực hiện thao tác này"})
		c.AbortWithStatusJSON(common.StatusForbidden, resp)
		// c.Next()
	}
}
