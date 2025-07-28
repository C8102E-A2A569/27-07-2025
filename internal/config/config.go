package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Limits LimitsConfig `yaml:"limits"`
	Files  FilesConfig  `yaml:"files"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LimitsConfig struct {
	MaxTasks        int `yaml:"max_tasks"`
	MaxFilesPerTask int `yaml:"max_files_per_task"`
}

type FilesConfig struct {
	AllowedExt      []string      `yaml:"allowed_extensions"`
	DownloadTimeout time.Duration `yaml:"download_timeout"`
	MaxSizeMB       int           `yaml:"max_size_mb"`
}

func MustLoad(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
