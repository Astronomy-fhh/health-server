package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type ServerConfig struct {
	ServerName string    `yaml:"server_name"`
	Env        string    `yaml:"env"`
	Gin        GinConfig `yaml:"gin"`
	Db         DbConfig  `yaml:"db"`
	S3         S3Config  `yaml:"s3"`
}

type DbConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Addr         string `yaml:"addr"`
	Dbname       string `yaml:"dbname"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type S3Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKey       string `yaml:"access_key"`
	SecretKey       string `yaml:"secret_key"`
	Region          string `yaml:"region"`
	ProductImageUri string `yaml:"product_image_uri"`
	UserAvatarUri   string `yaml:"user_avatar_uri"`
}

type GinConfig struct {
	Port string `yaml:"port"`
}

var serverConfigIst ServerConfig

func InitConfig(filePath string) error {
	var config ServerConfig
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	serverConfigIst = config
	return nil
}

func Get() ServerConfig {
	return serverConfigIst
}
