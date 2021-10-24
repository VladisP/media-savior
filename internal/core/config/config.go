package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/VladisP/media-savior/internal/common/validator"
)

const (
	configPathEnv = "CONFIG_PATH"
)

type Config struct {
	Env        string     `yaml:"env" validate:"required"`
	Logger     Logger     `yaml:"logger"`
	HTTPServer HTTPServer `yaml:"http_server" validate:"required"`
	DB         DB         `yaml:"db" validate:"required"`
}

type Logger struct {
	DevMode bool `yaml:"dev_mode"`
}

type HTTPServer struct {
	Port       string `yaml:"port" validate:"required"`
	Host       string `yaml:"host" validate:"required"`
	GinDevMode bool   `yaml:"gin_dev_mode"`
}

type DB struct {
	Host     string `yaml:"host" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	SSLMode  string `yaml:"ssl_mode" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
}

func readConfig(configEnvName string, config interface{}) error {
	configPath := os.Getenv(configEnvName)
	if configPath == "" {
		return fmt.Errorf("no config path: %s", configEnvName)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read config file: %s", configPath))
	}

	if err = yaml.Unmarshal(configBytes, config); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to unmarshal yaml config: %s", configPath))
	}
	return nil
}

func NewConfig(validator validator.Validator) (*Config, error) {
	var c Config

	if err := readConfig(configPathEnv, &c); err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	return &c, validator.ValidateStruct(c)
}
