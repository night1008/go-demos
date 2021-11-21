package app

import (
	"errors"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/model"
	"gin-api-server/pkg/auth"

	"gorm.io/gorm"
)

type PwdAuther struct {
	DB *gorm.DB
}

func (a *PwdAuther) GetDB() *gorm.DB {
	return a.DB
}

func (a *PwdAuther) Authenticate(email, password string) (bool, interface{}, error) {
	var u model.User
	if err := a.DB.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, svcerrorsx.ErrUserWrongEmailOrPassword
		}
	}

	if ok := auth.CheckPasswordHash(password, u.Password); !ok {
		return false, nil, svcerrorsx.ErrUserWrongEmailOrPassword
	} else {
		return true, &u, nil
	}
}
