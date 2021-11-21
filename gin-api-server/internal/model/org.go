package model

import (
	"gin-api-server/internal/extend/gormx"

	"github.com/jinzhu/copier"
)

type Org struct {
	gormx.Model
	Name      string `gorm:"not null"` // 组织名称
	Logo      string // 组织图标地址
	CreatorID uint64 // 创建用户ID
}

func (m Org) ToDTO() *OrgDTO {
	item := new(OrgDTO)
	copier.Copy(item, &m)
	return item
}

type Orgs []*Org

func (rs Orgs) ToDTOs() []*OrgDTO {
	list := make([]*OrgDTO, len(rs))
	for i, item := range rs {
		list[i] = item.ToDTO()
	}
	return list
}

type OrgDTO struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	CreatedAt int64  `json:"created_at"`
	IsOwner   bool   `json:"is_owner"`
}

type UserOrgDTOList []*OrgDTO

type OrgUserDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	Email        string `json:"email"`
	IsAdmin      bool   `json:"is_admin"`
	Confirmed    bool   `json:"confirmed"`
	LastActiveAt int64  `json:"last_active_at"`
	CreatedAt    int64  `json:"created_at"`
}
