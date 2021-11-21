package gormx

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Model struct {
	ID        uint64 `gorm:"primarykey"`           // 数据库ID
	CreatedAt int64  `gorm:"autoCreateTime:milli"` // 记录创建时间，单位 milli
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"` // 记录更新时间，单位 milli
}

type Config struct {
	Debug        bool
	Type         string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

// New Create gorm.DB instance
func New(c *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(c.Type) {
	case "mysql":
		// create database if not exists
		cfg, err := mysqlDriver.ParseDSN(c.DSN)
		if err != nil {
			return nil, err
		}

		err = createDatabaseWithMySQL(cfg)
		if err != nil {
			return nil, err
		}

		dialector = mysql.Open(c.DSN)
	case "postgres":
		dialector = postgres.Open(c.DSN)
	default:
		dialector = sqlite.Open(c.DSN)
	}

	gconfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: c.TablePrefix,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(dialector, gconfig)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	return db, nil
}

func createDatabaseWithMySQL(cfg *mysqlDriver.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", cfg.User, cfg.Passwd, cfg.Addr)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", cfg.DBName)
	_, err = db.Exec(query)
	return err
}

func FindOne(db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
