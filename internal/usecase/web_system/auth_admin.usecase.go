package web_system

import (
	"context"
	"errors"
	"time"

	websystemdto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/repository/model"
	websystemrepo "go-structure/internal/repository/web_system"
	websystemtransformer "go-structure/internal/transformer/web_system"
	"go-structure/internal/utils/generate"
	jwtutil "go-structure/internal/utils/jwt"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
)

type (
	IAuthAdminUsecase interface {
		LoginAdmin(ctx context.Context, req *websystemdto.AdminLoginRequestDto, ip string, userAgent string) (*websystemdto.AdminLoginResponseDto, error)
		RefreshToken(ctx context.Context, refreshToken string, ip string, userAgent string) (*websystemdto.AdminRefreshTokenResponseDto, error)
	}

	authAdminUsecase struct {
		adminRepo            websystemrepo.ISystemAdminRepository
		loginHistoryRepo     websystemrepo.ISystemLoginHistoryRepository
		refreshTokenRepo     websystemrepo.ISystemAdminRefreshTokenRepository
		txManager            database.TransactionManager
	}
)

var (
	ErrAdminNotFound         = errors.New("không tìm thấy admin")
	ErrAdminInvalidPassword  = errors.New("mật khẩu không chính xác")
	ErrAdminNotActive        = errors.New("tài khoản admin chưa được kích hoạt")
	ErrInvalidRefreshToken   = errors.New("refresh token không hợp lệ hoặc đã hết hạn")
)

const (
	refreshTokenTTL = 15 * 24 * time.Hour
)

func NewAuthAdminUsecase(
	adminRepo websystemrepo.ISystemAdminRepository,
	loginHistoryRepo websystemrepo.ISystemLoginHistoryRepository,
	refreshTokenRepo websystemrepo.ISystemAdminRefreshTokenRepository,
	txManager database.TransactionManager,
) IAuthAdminUsecase {
	return &authAdminUsecase{
		adminRepo:        adminRepo,
		loginHistoryRepo: loginHistoryRepo,
		refreshTokenRepo: refreshTokenRepo,
		txManager:        txManager,
	}
}

func (u *authAdminUsecase) LoginAdmin(ctx context.Context, req *websystemdto.AdminLoginRequestDto, ip string, userAgent string) (*websystemdto.AdminLoginResponseDto, error) {
	var (
		admin        *model.SystemAdmin
		refreshToken string
	)

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		admin, err = u.getAdminForLogin(txCtx, req.Email, req.Password)
		if err != nil {
			_ = u.logLoginHistory(txCtx, uuid.Nil, ip, userAgent, err, req.Email)
			return err
		}

		// Tạo refresh token và lưu vào database
		refreshToken = generate.GenerateRandomString(64)
		expiresAt := time.Now().Add(refreshTokenTTL)
		
		_, err = u.refreshTokenRepo.CreateRefreshToken(txCtx, &model.SystemAdminRefreshToken{
			AdminID:         admin.ID,
			RefreshTokenHash: refreshToken,
			ExpiresAt:       expiresAt,
			IpAddress:       ip,
			UserAgent:       userAgent,
		})
		if err != nil {
			return err
		}

		// Cập nhật last_login_at
		if err := u.adminRepo.UpdateLastLoginAt(txCtx, admin.ID); err != nil {
			return err
		}

		// Log login history thành công
		if err := u.logLoginHistory(txCtx, admin.ID, ip, userAgent, nil, ""); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Tạo access token với admin_id
	accessToken, err := jwtutil.GenerateAdminAccessToken(admin.ID)
	if err != nil {
		return nil, err
	}

	resp := websystemtransformer.ToAdminLoginResponseDto(accessToken, refreshToken, admin)
	return &resp, nil
}

func (u *authAdminUsecase) getAdminForLogin(ctx context.Context, email string, password string) (*model.SystemAdmin, error) {
	admin, err := u.adminRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	if ok := validator.CheckPassword(password, admin.PasswordHash); !ok {
		return nil, ErrAdminInvalidPassword
	}

	if !admin.IsActive {
		return nil, ErrAdminNotActive
	}

	return admin, nil
}

func (u *authAdminUsecase) logLoginHistory(
	ctx context.Context,
	adminID uuid.UUID,
	ip string,
	userAgent string,
	loginErr error,
	email string,
) error {
	var result string
	var reason string
	if loginErr == nil {
		result = "success"
	} else {
		if adminID == uuid.Nil && email != "" {
			admin, _ := u.adminRepo.GetByEmail(ctx, email)
			if admin != nil {
				adminID = admin.ID
			}
		}
		switch {
		case errors.Is(loginErr, ErrAdminNotFound):
			result = "failed_not_found"
			reason = "Tài khoản admin không tồn tại"
		case errors.Is(loginErr, ErrAdminInvalidPassword):
			result = "failed_password"
			reason = "Mật khẩu không chính xác"
		case errors.Is(loginErr, ErrAdminNotActive):
			result = "failed_inactive"
			reason = "Tài khoản admin chưa được kích hoạt"
		default:
			result = "failed_unknown"
			reason = loginErr.Error()
		}
	}
	_, err := u.loginHistoryRepo.CreateSystemLoginHistory(ctx, &model.SystemLoginHistory{
		AdminID:       adminID,
		Result:        result,
		FailureReason: reason,
		IpAddress:     ip,
		UserAgent:     userAgent,
	})
	return err
}

func (u *authAdminUsecase) RefreshToken(ctx context.Context, refreshToken string, ip string, userAgent string) (*websystemdto.AdminRefreshTokenResponseDto, error) {
	var (
		admin      *model.SystemAdmin
		newRefresh string
		newAccess  string
	)

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Lấy refresh token từ database
		tokenRecord, err := u.refreshTokenRepo.GetRefreshTokenByHash(txCtx, refreshToken)
		if err != nil || tokenRecord == nil {
			return ErrInvalidRefreshToken
		}

		// Kiểm tra token đã bị revoke hoặc hết hạn
		now := time.Now()
		if tokenRecord.IsRevoked || now.After(tokenRecord.ExpiresAt) {
			_ = u.refreshTokenRepo.RevokeRefreshTokenByHash(txCtx, refreshToken, "expired_or_revoked")
			return ErrInvalidRefreshToken
		}

		// Lấy thông tin admin
		admin, err = u.adminRepo.GetByID(txCtx, tokenRecord.AdminID)
		if err != nil || admin == nil {
			return ErrAdminNotFound
		}
		if !admin.IsActive {
			return ErrAdminNotActive
		}

		// Revoke token cũ
		if err := u.refreshTokenRepo.RevokeRefreshTokenByHash(txCtx, refreshToken, "rotated"); err != nil {
			return err
		}

		// Tạo refresh token mới
		newRefresh = generate.GenerateRandomString(64)
		newExpiresAt := now.Add(refreshTokenTTL)

		_, err = u.refreshTokenRepo.CreateRefreshToken(txCtx, &model.SystemAdminRefreshToken{
			AdminID:         admin.ID,
			RefreshTokenHash: newRefresh,
			ExpiresAt:       newExpiresAt,
			IpAddress:       ip,
			UserAgent:       userAgent,
		})
		if err != nil {
			return err
		}

		// Tạo access token mới với admin_id
		newAccess, err = jwtutil.GenerateAdminAccessToken(admin.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := websystemtransformer.ToAdminRefreshTokenResponseDto(newAccess, newRefresh)
	return &resp, nil
}
