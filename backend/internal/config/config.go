package config

import "github.com/spf13/viper"

type Config struct {
	Server            ServerConfig   `mapstructure:"server"`
	Database           DatabaseConfig `mapstructure:"database"`
	JWT                JWTConfig      `mapstructure:"jwt"`
	APIKey             string         `mapstructure:"api_key"`
	WinTillModulesPath string         `mapstructure:"win_till_modules_path"`
	LDAP               *LDAPConfig    `mapstructure:"ldap"`
	SSO                *SSOConfig     `mapstructure:"sso"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	TTL    int    `mapstructure:"ttl"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type LDAPConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	BaseDN   string `mapstructure:"base_dn"`
	BindDN   string `mapstructure:"bind_dn"`
	BindPass string `mapstructure:"bind_password"`
	Filter   string `mapstructure:"filter"`
	UIDAttr  string `mapstructure:"uid_attr"`
	MailAttr string `mapstructure:"mail_attr"`
	NameAttr string `mapstructure:"name_attr"`
	DefaultRoleID int64 `mapstructure:"default_role_id"`
	StartTLS bool   `mapstructure:"starttls"`
}

type SSOConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	ProviderURL  string `mapstructure:"provider_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	DefaultRoleID int64 `mapstructure:"default_role_id"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "cashier_config",
			SSLMode:  "disable",
		},
		JWT: JWTConfig{
			Secret: "change-me-in-production",
			TTL:    24,
		},
		WinTillModulesPath: "../win-till-modules",
	}
}
