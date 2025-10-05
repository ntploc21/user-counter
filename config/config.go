package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env          string          `mapstructure:"env"`
	IsProduction bool            `mapstructure:"isProduction"`
	Server       *ServerConfig   `mapstructure:"api"`
	Logger       *LoggerConfig   `mapstructure:"logger"`
	Database     *DatabaseConfig `mapstructure:"database"`
	Redis        *RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Host  string `mapstructure:"host"  default:"localhost"`
	Port  string `mapstructure:"port"  default:"8080"`
	Debug bool   `mapstructure:"debug" default:"true"`

	Swagger struct {
		Enabled bool `mapstructure:"enabled" default:"true"`
	} `mapstructure:"swagger"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"databaseName"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type RedisConfig struct {
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
}

var (
	_, b, _, _        = runtime.Caller(0)
	basePath          = filepath.Dir(b)
	defaultConfigFile = basePath + "/dev.yaml"
	v                 = viper.New()
	appConfig         AppConfig
)

func init() {
	Load()
}

func Load() {
	var configFile string
	if configFile = os.Getenv("CONFIG_PATH"); len(configFile) == 0 {
		configFile = defaultConfigFile
	}
	logrus.Infof("Loading config from %s", configFile)

	if err := loadConfigFile(configFile); err != nil {
		panic(err)
	}

	if err := scanConfigFile(&appConfig); err != nil {
		panic(err)
	}

	if err := validateConfig(&appConfig); err != nil {
		panic(err)
	}

	defaults.SetDefaults(&appConfig)
}

func loadConfigFile(configFile string) error {
	configFileName := filepath.Base(configFile)
	configFilePath := filepath.Dir(configFile)

	v.AddConfigPath(configFilePath)
	v.SetConfigName(
		strings.TrimSuffix(
			configFileName,
			filepath.Ext(configFileName),
		),
	)
	v.AutomaticEnv()

	return v.ReadInConfig()
}

func scanConfigFile(config any) error {
	return v.Unmarshal(&config)
}

func validateConfig(config any) error {
	validate := validator.New()
	return validate.Struct(config)
}

func GetAppConfig() *AppConfig {
	return &appConfig
}

func GetMySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appConfig.Database.Username,
		appConfig.Database.Password,
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.DatabaseName,
	)
}

func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", appConfig.Redis.Addrs[0], appConfig.Database.Port)
}
