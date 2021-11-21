package model

import (
	"gin-api-server/internal/extend/gormx"

	"github.com/jinzhu/copier"
)

const (
	PermissionUpdateOrg                  = "org:update"
	PermissionDeleteOrg                  = "org:delete"
	PermissionListOrgDeletedApp          = "org:list-deleted-app"
	PermissionRestoreOrgDeletedApp       = "org:restore-deleted-app"
	PermissionGetApp                     = "app:get"
	PermissionCreateApp                  = "app:create"
	PermissionUpdateApp                  = "app:update"
	PermissionDeleteApp                  = "app:delete"
	PermissionUpdateAppMember            = "app:update-member"
	PermissionInviteAppMember            = "app:invite-member"
	PermissionResendAppMemberInviteEmail = "app:resend-member-invite-email"
)

type Permission struct {
	gormx.Model
	Name string `gorm:"not null;uniqueIndex"`
}

func (p Permission) ToDTO() *PermissionDTO {
	item := new(PermissionDTO)
	copier.Copy(item, &p)
	return item
}

type PermissionDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Permissions []*Permission

func (rs Permissions) ToDTOs() []*PermissionDTO {
	list := make([]*PermissionDTO, len(rs))
	for i, item := range rs {
		list[i] = item.ToDTO()
	}
	return list
}
