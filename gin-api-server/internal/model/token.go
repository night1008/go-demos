package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

const TokenTypeRegister = "register"

type Token struct {
	ID        string `gorm:"primarykey"`           // token id, 当前为 uuid 形式
	Email     string `gorm:"not null"`             // 用户注册邮箱
	Type      string `gorm:"not null"`             // token 类型
	ExpiredAt int64  `gorm:"not null"`             // token 过期时间，单位 milli
	CreatedAt int64  `gorm:"autoCreateTime:milli"` // token 创建时间，单位 milli
}

func (u *Token) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

func (m Token) ToDTO() *TokenDTO {
	item := new(TokenDTO)
	copier.Copy(item, &m)
	return item
}

type TokenDTO struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Type            string    `json:"type"`
	ExpiredAt       int64     `json:"expired_at"`
	CreatedAt       int64     `json:"created_at"`
	EmailVerifiedAt int64     `json:"email_verify_at"`
}
