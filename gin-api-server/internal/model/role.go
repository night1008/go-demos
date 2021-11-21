package model

import (
	"gin-api-server/internal/extend/gormx"

	"github.com/jinzhu/copier"
)

const (
	RoleOrgOwner  = "组织拥有者"
	RoleAppAdmin  = "应用管理员"
	RoleAppMember = "应用普通成员"
)

type Role struct {
	gormx.Model
	Name string `gorm:"not null;uniqueIndex"` // 角色名称
}

func (r Role) ToDTO() *RoleDTO {
	item := new(RoleDTO)
	copier.Copy(item, &r)
	return item
}

type Roles []*Role

func (rs Roles) ToDTOs() []*RoleDTO {
	list := make([]*RoleDTO, len(rs))
	for i, item := range rs {
		list[i] = item.ToDTO()
	}
	return list
}

type RoleDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type RoleDTOList []*RoleDTO

type RoleDTOPageList struct {
	Total    int64      `json:"total"`
	Current  int        `json:"current"`
	PageSize int        `json:"page_size"`
	List     []*RoleDTO `json:"list"`
}
