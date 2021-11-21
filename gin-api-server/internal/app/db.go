package app

import (
	"errors"
	"gin-api-server/internal/config"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/model"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type DB struct {
	Cfg *config.DatabaseCfg
}

func (d *DB) Setup() (*gorm.DB, error) {
	switch d.Cfg.Type {
	case "mysql", "postgres":
	case "sqlite3":
		_ = os.MkdirAll(filepath.Dir(d.Cfg.DSN), 0777)
	default:
		return nil, errors.New("unknown db")
	}

	db, err := gormx.New(&gormx.Config{
		Debug:        d.Cfg.Debug,
		Type:         d.Cfg.Type,
		DSN:          d.Cfg.DSN,
		MaxIdleConns: d.Cfg.MaxIdleConns,
		MaxLifetime:  d.Cfg.MaxLifetime,
		MaxOpenConns: d.Cfg.MaxOpenConns,
		TablePrefix:  d.Cfg.TablePrefix,
	})
	if err != nil {
		return nil, err
	}

	if d.Cfg.EnableAutoMigrate {
		err = AutoMigrate(db, d.Cfg.Type)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// MEMO 后续新增 model，需要修改此处代码
func AutoMigrate(db *gorm.DB, dbType string) error {
	if strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(model.User),
		new(model.Org),
		new(model.App),
		new(model.Role),
		new(model.Permission),
		new(model.RolePermission),
		new(model.UserRole),
		new(model.Token),
	)
}
