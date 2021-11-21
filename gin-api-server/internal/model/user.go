package model

import (
	"gin-api-server/internal/extend/gormx"

	"github.com/jinzhu/copier"
)

type User struct {
	gormx.Model
	Name            string `gorm:"not null;index"` // 用户名
	Logo            string // 用户图标地址
	Email           string `gorm:"not null;uniqueIndex"` // 用户邮箱
	Password        string // 加密后的密码
	IsAdmin         bool   `gorm:"default:false"` // 是否平台管理员
	EmailVerifiedAt int64  // 用户邮箱验证时间，单位 milli
	LastActiveAt    int64  // 用户最后活跃时间，单位 milli
}

func (m User) ToDTO() *UserDTO {
	item := new(UserDTO)
	copier.Copy(item, &m)
	return item
}

type Users []*User

func (us Users) ToDTOs() []*UserDTO {
	list := make([]*UserDTO, len(us))
	for i, item := range us {
		list[i] = item.ToDTO()
	}
	return list
}

type UserDTO struct {
	ID              uint64 `json:"id" example:"1"`
	Name            string `json:"name" example:"demo"`
	Logo            string `json:"logo" example:" "`
	Password        string `json:"-"`
	Email           string `json:"email" example:"demo@xmfunny.com"`
	IsAdmin         bool   `json:"is_admin" example:"false"`
	EmailVerifiedAt int64  `json:"email_verified_at" example:"1635696000000"`
	LastActiveAt    int64  `json:"last_active_at" example:"0"`
}

type UserDTOList []*UserDTO

type UserAppDTO struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Logo      string `json:"logo"`
	RoleID    uint64 `json:"role_id"`
	CreatedAt int64  `json:"created_at"`
}

type UserWithAppDTO struct {
	UserDTO
	IsOwner bool      `json:"is_owner"`
	Apps    []*AppDTO `json:"apps"`
}

type UserWithAppDTOList []*UserWithAppDTO
