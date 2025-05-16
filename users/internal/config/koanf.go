package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type AppConfig struct {
	ConfigService *koanf.Koanf
}

func NewAppConfig(k *koanf.Koanf, path string) (*AppConfig, error) {
	config := &AppConfig{
		ConfigService: k,
	}
	err := config.load(path)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (appConfig *AppConfig) load(path string) error {
	if err := appConfig.ConfigService.Load(file.Provider(path), yaml.Parser()); err != nil {
		return err
	}
	fmt.Println("APP PORT ->>>>", appConfig.ConfigService.String("app.port"))
	fmt.Println("gRPC PORT ->>>>", appConfig.ConfigService.String("app.grpc_port"))
	return nil
}
