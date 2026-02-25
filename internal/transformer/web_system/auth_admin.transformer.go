package web_system

import (
	"time"

	websystemdto "go-structure/internal/dto/web_system"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
)

func GetUniquePermissionIDs(permMap map[uuid.UUID][]uuid.UUID) []uuid.UUID {
	if len(permMap) == 0 {
		return nil
	}
	set := make(map[uuid.UUID]struct{})
	for _, ids := range permMap {
		for _, id := range ids {
			set[id] = struct{}{}
		}
	}
	out := make([]uuid.UUID, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func ToAdminRolesAndPermissionsDtos(
	roles []*websystem.Role,
	permMap map[uuid.UUID][]uuid.UUID,
	permissions []*websystem.Permission,
) ([]websystemdto.RoleItemDto, []websystemdto.PermissionItemDto) {
	roleDtos := make([]websystemdto.RoleItemDto, 0, len(roles))
	for _, r := range roles {
		permIDs := permMap[r.ID]
		permIDStrs := permissionUUIDsToStrings(permIDs)
		roleDtos = append(roleDtos, ToRoleItemDto(r, permIDStrs))
	}
	permDtos := make([]websystemdto.PermissionItemDto, 0, len(permissions))
	for _, p := range permissions {
		permDtos = append(permDtos, ToPermissionItemDto(p))
	}
	return roleDtos, permDtos
}

func ToAdminLoginResponseDto(
	accessToken string,
	refreshToken string,
	accessTokenExpiresAt time.Time,
	admin *model.SystemAdmin,
	roles []websystemdto.RoleItemDto,
	permissions []websystemdto.PermissionItemDto,
) websystemdto.AdminLoginResponseDto {
	roleSimples := toAdminRoleSimpleDtos(roles)
	permSimples := toAdminPermissionSimpleDtos(permissions)
	now := time.Now()
	expiresIn := int64(0)
	if accessTokenExpiresAt.After(now) {
		expiresIn = int64(accessTokenExpiresAt.Sub(now).Seconds())
	}
	return websystemdto.AdminLoginResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		ExpiredAt:    accessTokenExpiresAt.UTC().Format(time.RFC3339),
		Admin: websystemdto.AdminProfileResponseDto{
			ID:          admin.ID,
			Email:       admin.Email,
			FullName:    admin.FullName,
			Department:  admin.Department,
			IsActive:    admin.IsActive,
			Roles:       roleSimples,
			Permissions: permSimples,
		},
	}
}

func toAdminRoleSimpleDtos(roles []websystemdto.RoleItemDto) []websystemdto.AdminRoleSimpleDto {
	if len(roles) == 0 {
		return []websystemdto.AdminRoleSimpleDto{}
	}
	out := make([]websystemdto.AdminRoleSimpleDto, 0, len(roles))
	for _, r := range roles {
		out = append(out, websystemdto.AdminRoleSimpleDto{ID: r.ID, Name: r.Name})
	}
	return out
}

func toAdminPermissionSimpleDtos(permissions []websystemdto.PermissionItemDto) []websystemdto.AdminPermissionSimpleDto {
	if len(permissions) == 0 {
		return []websystemdto.AdminPermissionSimpleDto{}
	}
	out := make([]websystemdto.AdminPermissionSimpleDto, 0, len(permissions))
	for _, p := range permissions {
		out = append(out, websystemdto.AdminPermissionSimpleDto{ID: p.ID, Code: p.Code, Name: p.Name})
	}
	return out
}

func ToAdminRefreshTokenResponseDto(accessToken string, refreshToken string) websystemdto.AdminRefreshTokenResponseDto {
	return websystemdto.AdminRefreshTokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
