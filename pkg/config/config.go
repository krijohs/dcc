package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ConfigPath string `mapstructure:"config_path"`
	KubeConf   string `mapstructure:"kubeconf"`
	Registries []DockerRegistry
	LogFormat  string `mapstructure:"log_format"`
	LogLevel   string `mapstructure:"log_level"`
	LogFile    string `mapstructure:"log_file"`
}

type DockerRegistry struct {
	Name    string
	Config  string
	Exclude []string
	Include []string
}

// Load config using config.yaml file, values can be overridden by environment variables.
func Load() (Config, error) {
	cfg := Config{}

	viper.SetDefault("log_format", "text")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_file", "stdout")
	viper.SetEnvPrefix("dcc")
	viper.AutomaticEnv()

	viper.SetConfigName("config")

	cfgPath := viper.GetString("config_path")
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dockerconfig")
	viper.AddConfigPath("/etc/dockerconfig/")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, errors.Wrap(err, "unable to read config")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, errors.Wrap(err, "unable to unmarshal config")
	}

	return cfg, nil
}
