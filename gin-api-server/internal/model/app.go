package model

import (
	"gin-api-server/internal/extend/gormx"

	"github.com/jinzhu/copier"
	"gorm.io/plugin/soft_delete"
)

const (
	AppDeployStatusWaiting   = "waiting"
	AppDeployStatusDeploying = "deploying"
	AppDeployStatusSuccess   = "success"
	AppDeployStatusFailed    = "failed"
)

const (
	AppStatus = ""
)

type App struct {
	gormx.Model
	DeletedAt       soft_delete.DeletedAt `gorm:"softDelete:milli;index;"` // 软删除时间，单位 milli
	Name            string                `gorm:"not null"`                // 应用名称
	Logo            string                // 应用图标地址
	Identifier      string                `gorm:"not null;uniqueIndex"` // 应用英文标识符
	AccessKeyID     string                `gorm:"not null;uniqueIndex"` // 应用访问Key
	AccessKeySecret string                `gorm:"not null"`             // 应用访问密钥
	DeployStatus    string                // 应用底层资源部署状态
	Status          string                // 应用状态

	OrgID     uint64 `gorm:"index"` // 组织ID
	CreatorID uint64 // 创建用户ID
	DeletorID uint64 // 删除用户ID
}

func (m App) ToDTO() *AppDTO {
	item := new(AppDTO)
	copier.Copy(item, &m)
	return item
}

type Apps []*App

func (as Apps) ToDTOs() []*AppDTO {
	list := make([]*AppDTO, len(as))
	for i, item := range as {
		list[i] = item.ToDTO()
	}
	return list
}

type AppDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Identifier  string `json:"identifier"`
	AccessKeyID string `json:"access_key_id"`
	CreatedAt   int64  `json:"created_at"`
	OrgID       uint64 `json:"org_id"`
	RoleID      uint64 `json:"role_id"`
}

type AppDTOList []*AppDTO

type DeletedAppDTO struct {
	AppDTO
	DeletedAt int64    `json:"deleted_at"`
	DeletorID uint64   `json:"deletor_id"`
	DeletedBy *UserDTO `json:"deleted_by"`
}

func (m App) ToDeletedAppDTO() *DeletedAppDTO {
	item := new(DeletedAppDTO)
	copier.Copy(item, &m)
	return item
}

type DeletedAppDTOList []*DeletedAppDTO

func (as Apps) ToDeletedAppDTOs() []*DeletedAppDTO {
	list := make([]*DeletedAppDTO, len(as))
	for i, item := range as {
		list[i] = item.ToDeletedAppDTO()
	}
	return list
}
