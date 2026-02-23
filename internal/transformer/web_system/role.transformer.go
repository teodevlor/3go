package web_system

import (
	"go-structure/internal/dto/common"
	dto_websystem "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
)

func ToRoleItemDto(r *websystem.Role, permissionIDs []string) dto_websystem.RoleItemDto {
	if r == nil {
		return dto_websystem.RoleItemDto{}
	}
	if permissionIDs == nil {
		permissionIDs = []string{}
	}
	return dto_websystem.RoleItemDto{
		ID:            r.ID.String(),
		Code:          r.Code,
		Name:          r.Name,
		Description:   r.Description,
		IsActive:      r.IsActive,
		PermissionIDs: permissionIDs,
	}
}

func ToListRolesResponseDto(roles []*websystem.Role, permMap map[uuid.UUID][]uuid.UUID, page, limit int, total int64) dto_websystem.ListRolesResponseDto {
	items := make([]dto_websystem.RoleItemDto, 0, len(roles))
	for _, r := range roles {
		permIDStrs := permissionUUIDsToStrings(permMap[r.ID])
		items = append(items, ToRoleItemDto(r, permIDStrs))
	}
	return dto_websystem.ListRolesResponseDto{
		Items: items,
		Pagination: common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}

func permissionUUIDsToStrings(uuids []uuid.UUID) []string {
	if len(uuids) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(uuids))
	for _, u := range uuids {
		out = append(out, u.String())
	}
	return out
}
