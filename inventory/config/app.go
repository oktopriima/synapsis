package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Port   string `mapstructure:"port"`
		Domain string `mapstructure:"domain"`
		Name   string `mapstructure:"name"`
	} `mapstructure:"app"`
	Postgres struct {
		Host               string `mapstructure:"host"`
		Database           string `mapstructure:"database"`
		Password           string `mapstructure:"password"`
		Port               string `mapstructure:"port"`
		User               string `mapstructure:"user"`
		MigrationDirectory string `mapstructure:"migration_directory"`
		LogDirectory       string `mapstructure:"log_directory"`
		LogFile            string `mapstructure:"log_file"`
		Zone               string `mapstructure:"zone"`
	} `mapstructure:"postgres"`
}

func NewAppConfig() (app AppConfig) {
	path := os.Getenv("CONFIG_PATH")
	osEnv := os.Getenv("OS_ENV")

	env := "env"
	if osEnv != "" {
		env = osEnv
	}

	envFile := fmt.Sprintf("%s.yaml", env)
	if path != "" {
		envFile = path + "/" + env + ".yaml"
	}

	replacer := strings.NewReplacer(`.`, `_`)
	viper.AddConfigPath(path)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(`yaml`)
	viper.SetConfigFile(envFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&app)
	if err != nil {
		panic(err)
	}

	return
}
