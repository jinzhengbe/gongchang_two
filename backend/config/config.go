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
		BaseURL        string   `yaml:"base_url"`
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
	
	// 处理数据库连接环境变量
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}
	
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