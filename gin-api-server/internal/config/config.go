package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	once sync.Once
)

func Load(path string) (cfg *Config, err error) {
	once.Do(func() {
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			return
		}

		if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
			return
		}
	})
	return cfg, err
}

type Config struct {
	Logger     LoggerCfg     `yaml:"logger"`
	HTTPServer HTTPServerCfg `yaml:"http"`
	Database   DatabaseCfg   `yaml:"database"`
	Seeder     SeederCfg     `yaml:"seeder"`
}

type LoggerCfg struct {
	Level string `yaml:"level"`
}

type HTTPServerCfg struct {
	RunMode            string         `yaml:"run_mode"`
	Addr               string         `yaml:"addr"`
	CertFile           string         `yaml:"cert_file"`
	KeyFile            string         `yaml:"key_file"`
	ShutdownTimeout    int64          `yaml:"shutdown_timeout"`
	MaxContentLength   uint           `yaml:"max_content_length"`
	MaxReqLoggerLength uint           `yaml:"max_req_logger_length"`
	CORS               CORSCfg        `yaml:"cors"`
	Session            HTTPSessionCfg `yaml:"session"`
	Swagger            SwaggerCfg     `yaml:"swagger"`
}

type HTTPSessionCfg struct {
	Name           string `yaml:"name"`
	Secret         string `yaml:"secret"`
	ExpireDuration uint   `yaml:"expire_duration"`
	Path           string `yaml:"path"`
}

type CORSCfg struct {
	Enable           bool     `yaml:"enable"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_hredentials"`
	MaxAge           uint     `yaml:"max_age"`
}

type SwaggerCfg struct {
	Enable bool `yaml:"enable"`
}

type DatabaseCfg struct {
	Debug             bool   `yaml:"debug"`
	Type              string `yaml:"type"`
	DSN               string `yaml:"dsn"`
	MaxLifetime       int    `yaml:"max_lifetime"`
	MaxOpenConns      int    `yaml:"max_open_conns"`
	MaxIdleConns      int    `yaml:"max_idle_conns"`
	TablePrefix       string `yaml:"table_prefix"`
	EnableAutoMigrate bool   `yaml:"enable_auto_migrate"`
}

type SeederCfg struct {
	Users []SeederUserCfg `yaml:"users"`
	Roles []SeederRoleCfg `yaml:"roles"`
}

type SeederUserCfg struct {
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	IsAdmin  bool   `yaml:"is_admin"`
}

type SeederRoleCfg struct {
	Name        string   `yaml:"name"`
	Permissions []string `yaml:"permissions"`
}
