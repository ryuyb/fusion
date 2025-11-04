package config

import "time"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Jwt      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	//Driver          string        `mapstructure:"driver"`
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type LoggerConfig struct {
	Level         string `mapstructure:"level"` // debug, info, warn, error
	EnableConsole bool   `mapstructure:"enable_console"`
	EnableColor   bool   `mapstructure:"enable_color"`
	EnableFile    bool   `mapstructure:"enable_file"`
	OutputPath    string `mapstructure:"output_path"`
	Compress      bool   `mapstructure:"compress"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxAge        int    `mapstructure:"max_age"`
	MaxBackups    int    `mapstructure:"max_backups"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}
