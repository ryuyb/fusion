package config

import "time"

type Config struct {
	App      AppConfig            `mapstructure:"app"`
	Server   ServerConfig         `mapstructure:"server"`
	Database DatabaseConfig       `mapstructure:"database"`
	Log      LogConfig            `mapstructure:"logger"`
	JWT      JWTConfig            `mapstructure:"jwt"`
	Job      map[string]JobConfig `mapstructure:"job"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env" validate:"oneof=dev prod"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port" validate:"min=1,max=65535"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port" validate:"min=1,max=65535"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SSLMode     string `mapstructure:"ssl_mode" validate:"oneof=disable require verify-ca verify-full"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

type LogConfig struct {
	Level         string `mapstructure:"level" validate:"oneof=debug info warn error"`
	EnableConsole bool   `mapstructure:"enable_console"`
	EnableFile    bool   `mapstructure:"enable_file"`
	ConsoleColor  bool   `mapstructure:"console_color"`
	FilePath      string `mapstructure:"file_path"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type JobConfig struct {
	Enable   bool   `mapstructure:"enable"`
	CronExpr string `mapstructure:"cron_expr"`
}
