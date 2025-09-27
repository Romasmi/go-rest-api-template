package config

import (
	"reflect"
	"strings"
	"time"

	"github.com/Romasmi/go-rest-api-template/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"` // mapping in annotation is optional and by default is use property name as it is
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         uint
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxConnections  uint
	MinConnections  uint
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret        string
	ExpirationTTL time.Duration
}

func bindEnvRecursive(v *viper.Viper, prefix string, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = utils.FirstChatToLowerCase(field.Name)
		}

		fieldPath := prefix
		if prefix != "" {
			fieldPath = prefix + "." + tag
		} else {
			fieldPath = tag
		}

		if field.Type.Kind() == reflect.Struct {
			if err := bindEnvRecursive(v, fieldPath, val.Field(i)); err != nil {
				return err
			}
			continue
		}

		envVar := strings.ToUpper(strings.ReplaceAll(fieldPath, ".", "_"))
		if err := v.BindEnv(fieldPath, envVar); err != nil {
			return err
		}
	}
	return nil
}

func bindAllEnvVars(v *viper.Viper) error {
	return bindEnvRecursive(v, "", reflect.ValueOf(&Config{}).Elem())
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v2 := viper.New()
	v2.SetConfigName("override")
	v2.SetConfigType("yaml")
	v2.AddConfigPath(configPath)
	if err := v2.ReadInConfig(); err == nil {
		err := v.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return nil, err
		}
	}

	if err := bindAllEnvVars(v); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
