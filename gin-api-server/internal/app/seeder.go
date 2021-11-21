package app

import (
	"gin-api-server/internal/authority"
	"gin-api-server/internal/config"
	"gin-api-server/internal/extend/timex"
	"gin-api-server/internal/model"
	"gin-api-server/internal/service"
	"gin-api-server/pkg/auth"

	"gorm.io/gorm"
)

type Seeder struct {
	Cfg *config.SeederCfg
	DB  *gorm.DB
}

func (s *Seeder) Setup() error {
	if err := initUsers(s.DB, s.Cfg.Users); err != nil {
		return err
	}
	if err := initRoles(s.DB, s.Cfg.Roles); err != nil {
		return err
	}
	return nil
}

func initUsers(db *gorm.DB, cfgUsers []config.SeederUserCfg) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, u := range cfgUsers {
			hashPassword, err := auth.HashPassword(u.Password, service.UserPasswordHashCost)
			if err != nil {
				return err
			}
			user := model.User{
				Name:            u.Name,
				Email:           u.Email,
				Password:        hashPassword,
				IsAdmin:         u.IsAdmin,
				EmailVerifiedAt: timex.GetNowUTCMilli(),
			}
			if err := tx.FirstOrCreate(&user, model.User{Email: u.Email}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func initRoles(db *gorm.DB, cfgRoles []config.SeederRoleCfg) error {
	// TODO 执行前先清空权限相关表
	a := authority.Authority{DB: db}
	for _, r := range cfgRoles {
		if _, err := a.CreateRole(r.Name); err != nil {
			return err
		}
		for _, p := range r.Permissions {
			if _, err := a.CreatePermission(p); err != nil {
				return err
			}
		}
		if err := a.AssignPermissions(r.Name, r.Permissions); err != nil {
			return err
		}
	}
	return nil
}
