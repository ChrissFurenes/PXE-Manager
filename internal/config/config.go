package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerIP string         `yaml:"server_ip"`
	TFTP     TFTPConfig     `yaml:"tftp"`
	HTTP     HTTPConfig     `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
}

type TFTPConfig struct {
	ListenAddr string `yaml:"listen_addr"`
	RootDir    string `yaml:"root_dir"`
}

type HTTPConfig struct {
	ListenAddr  string `yaml:"listen_addr"`
	RootDir     string `yaml:"root_dir"`
	TemplateDir string `yaml:"template_dir"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
