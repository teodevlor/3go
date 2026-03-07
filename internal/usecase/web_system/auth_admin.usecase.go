package web_system

import (
	"context"
	"errors"
	"time"

	"go-structure/global"
	websystemdto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	"go-structure/internal/repository/model"
	websystemrepo "go-structure/internal/repository/web_system"
	websystemtransformer "go-structure/internal/transformer/web_system"
	"go-structure/internal/utils/generate"
	jwtutil "go-structure/internal/utils/jwt"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IAuthAdminUsecase interface {
		LoginAdmin(ctx context.Context, req *websystemdto.AdminLoginRequestDto, ip string, userAgent string) (*websystemdto.AdminLoginResponseDto, error)
		RefreshToken(ctx context.Context, refreshToken string, ip string, userAgent string) (*websystemdto.AdminRefreshTokenResponseDto, error)
		AdminHasPermission(ctx context.Context, adminID uuid.UUID, permissionCode string) (bool, error)
	}

	authAdminUsecase struct {
		adminRepo          websystemrepo.ISystemAdminRepository
		loginHistoryRepo   websystemrepo.ISystemLoginHistoryRepository
		refreshTokenRepo   websystemrepo.ISystemAdminRefreshTokenRepository
		adminRoleRepo      websystemrepo.ISystemAdminRoleRepository
		roleRepo           websystemrepo.IRoleRepository
		rolePermissionRepo websystemrepo.IRolePermissionRepository
		permissionRepo     websystemrepo.IPermissionRepository
		txManager          database.TransactionManager
	}
)

var (
	ErrAdminNotFound        = errors.New("không tìm thấy admin")
	ErrAdminInvalidPassword = errors.New("mật khẩu không chính xác")
	ErrAdminNotActive       = errors.New("tài khoản admin chưa được kích hoạt")
	ErrInvalidRefreshToken  = errors.New("refresh token không hợp lệ hoặc đã hết hạn")
)

const (
	refreshTokenTTL = 15 * 24 * time.Hour
)

func NewAuthAdminUsecase(
	adminRepo websystemrepo.ISystemAdminRepository,
	loginHistoryRepo websystemrepo.ISystemLoginHistoryRepository,
	refreshTokenRepo websystemrepo.ISystemAdminRefreshTokenRepository,
	adminRoleRepo websystemrepo.ISystemAdminRoleRepository,
	roleRepo websystemrepo.IRoleRepository,
	rolePermissionRepo websystemrepo.IRolePermissionRepository,
	permissionRepo websystemrepo.IPermissionRepository,
	txManager database.TransactionManager,
) IAuthAdminUsecase {
	return &authAdminUsecase{
		adminRepo:          adminRepo,
		loginHistoryRepo:   loginHistoryRepo,
		refreshTokenRepo:   refreshTokenRepo,
		adminRoleRepo:      adminRoleRepo,
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		permissionRepo:     permissionRepo,
		txManager:          txManager,
	}
}

func (u *authAdminUsecase) LoginAdmin(ctx context.Context, req *websystemdto.AdminLoginRequestDto, ip string, userAgent string) (*websystemdto.AdminLoginResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("LoginAdmin: start", zap.String(global.KeyCorrelationID, cid), zap.String("email", req.Email))

	type txResult struct {
		admin        *model.SystemAdmin
		refreshToken string
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			admin, err := u.getAdminForLogin(txCtx, req.Email, req.Password)
			if err != nil {
				global.Logger.Error("LoginAdmin: failed to get admin for login", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
				_ = u.logLoginHistory(txCtx, uuid.Nil, ip, userAgent, err, req.Email)
				return txResult{}, err
			}

			refreshToken := generate.GenerateRandomString(64)
			expiresAt := time.Now().Add(refreshTokenTTL)

			_, err = u.refreshTokenRepo.CreateRefreshToken(txCtx, &model.SystemAdminRefreshToken{
				AdminID:          admin.ID,
				RefreshTokenHash: refreshToken,
				ExpiresAt:        expiresAt,
				IpAddress:        ip,
				UserAgent:        userAgent,
			})
			if err != nil {
				return txResult{}, err
			}

			if err := u.adminRepo.UpdateLastLoginAt(txCtx, admin.ID); err != nil {
				return txResult{}, err
			}

			if err := u.logLoginHistory(txCtx, admin.ID, ip, userAgent, nil, ""); err != nil {
				return txResult{}, err
			}

			return txResult{admin: admin, refreshToken: refreshToken}, nil
		},
	)
	if err != nil {
		global.Logger.Error("LoginAdmin: failed", zap.String(global.KeyCorrelationID, cid), zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	accessToken, accessTokenExpiresAt, err := jwtutil.GenerateAdminAccessToken(res.admin.ID)
	if err != nil {
		global.Logger.Error("LoginAdmin: failed to generate access token", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	global.Logger.Info("LoginAdmin: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("email", req.Email), zap.String("admin_id", res.admin.ID.String()))
	roles, permissions := u.getAdminRolesAndPermissions(ctx, res.admin.ID)
	resp := websystemtransformer.ToAdminLoginResponseDto(accessToken, res.refreshToken, accessTokenExpiresAt, res.admin, roles, permissions)
	return &resp, nil
}

func (u *authAdminUsecase) AdminHasPermission(ctx context.Context, adminID uuid.UUID, permissionCode string) (bool, error) {
	if permissionCode == "" {
		return false, nil
	}
	_, permissions := u.getAdminRolesAndPermissions(ctx, adminID)
	for _, p := range permissions {
		if p.Code == permissionCode {
			return true, nil
		}
	}
	return false, nil
}

func (u *authAdminUsecase) getAdminRolesAndPermissions(ctx context.Context, adminID uuid.UUID) ([]websystemdto.RoleItemDto, []websystemdto.PermissionItemDto) {
	if u.adminRoleRepo == nil || u.roleRepo == nil || u.rolePermissionRepo == nil || u.permissionRepo == nil {
		return nil, nil
	}
	roleIDs, err := u.adminRoleRepo.GetRoleIDsByAdminID(ctx, adminID)
	if err != nil || len(roleIDs) == 0 {
		return nil, nil
	}
	roles, err := u.roleRepo.GetByIDs(ctx, roleIDs)
	if err != nil || len(roles) == 0 {
		return nil, nil
	}
	permMap, err := u.rolePermissionRepo.GetPermissionIDsByRoleIDs(ctx, roleIDs)
	if err != nil {
		permMap = nil
	}
	permIDs := websystemtransformer.GetUniquePermissionIDs(permMap)
	if len(permIDs) == 0 {
		return websystemtransformer.ToAdminRolesAndPermissionsDtos(roles, permMap, nil)
	}
	permissions, err := u.permissionRepo.GetByIDs(ctx, permIDs)
	if err != nil || len(permissions) == 0 {
		return websystemtransformer.ToAdminRolesAndPermissionsDtos(roles, permMap, nil)
	}
	return websystemtransformer.ToAdminRolesAndPermissionsDtos(roles, permMap, permissions)
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("RefreshToken: start", zap.String(global.KeyCorrelationID, cid))

	type txResult struct {
		admin      *model.SystemAdmin
		newRefresh string
		newAccess  string
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			tokenRecord, err := u.refreshTokenRepo.GetRefreshTokenByHash(txCtx, refreshToken)
			if err != nil || tokenRecord == nil {
				return txResult{}, ErrInvalidRefreshToken
			}

			now := time.Now()
			if tokenRecord.IsRevoked || now.After(tokenRecord.ExpiresAt) {
				_ = u.refreshTokenRepo.RevokeRefreshTokenByHash(txCtx, refreshToken, "expired_or_revoked")
				return txResult{}, ErrInvalidRefreshToken
			}

			admin, err := u.adminRepo.GetByID(txCtx, tokenRecord.AdminID)
			if err != nil || admin == nil {
				return txResult{}, ErrAdminNotFound
			}
			if !admin.IsActive {
				return txResult{}, ErrAdminNotActive
			}

			if err := u.refreshTokenRepo.RevokeRefreshTokenByHash(txCtx, refreshToken, "rotated"); err != nil {
				return txResult{}, err
			}

			newRefresh := generate.GenerateRandomString(64)
			_, err = u.refreshTokenRepo.CreateRefreshToken(txCtx, &model.SystemAdminRefreshToken{
				AdminID:          admin.ID,
				RefreshTokenHash: newRefresh,
				ExpiresAt:        now.Add(refreshTokenTTL),
				IpAddress:        ip,
				UserAgent:        userAgent,
			})
			if err != nil {
				return txResult{}, err
			}

			newAccess, _, err := jwtutil.GenerateAdminAccessToken(admin.ID)
			if err != nil {
				return txResult{}, err
			}

			return txResult{admin: admin, newRefresh: newRefresh, newAccess: newAccess}, nil
		},
	)
	if err != nil {
		global.Logger.Error("RefreshToken: failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	global.Logger.Info("RefreshToken: completed successfully", zap.String(global.KeyCorrelationID, cid))
	resp := websystemtransformer.ToAdminRefreshTokenResponseDto(res.newAccess, res.newRefresh)
	return &resp, nil
}
