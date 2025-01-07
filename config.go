package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		BotToken string `yaml:"bot_token"`
		Debug    bool   `yaml:"debug"`
		Timeout  int    `yaml:"timeout"`
	} `yaml:"telegram"`
	System struct {
		RetryInterval int `yaml:"retry_interval"`
	} `yaml:"system"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
