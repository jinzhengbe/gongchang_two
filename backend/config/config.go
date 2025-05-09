package config

import (
	"os"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type Config struct {
	Server struct {
		Host           string   `yaml:"host"`
		Port           string   `yaml:"port"`
		TrustedProxies []string `yaml:"trusted_proxies"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	JWT struct {
		Secret string `yaml:"secret"`
		Expire int    `yaml:"expire"`
	} `yaml:"jwt"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	
	// 读取配置文件
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}
	
	// 解析 YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// 处理环境变量
	config.JWT.Secret = getEnvValue(config.JWT.Secret)
	
	return config, nil
}

func getEnvValue(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
		if envValue, exists := os.LookupEnv(envVar); exists {
			return envValue
		}
	}
	return value
} 