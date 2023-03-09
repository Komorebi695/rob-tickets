package configs

import (
	"gopkg.in/yaml.v2"
	"os"
)

// AppConfig 服务端配置
type AppConfig struct {
	AppName    string   `yaml:"app_name"`
	Port       string   `yaml:"port"`
	StaticPath string   `yaml:"static_path"`
	Mode       string   `yaml:"mode"`
	DataBase   DataBase `yaml:"data_base"`
}

// DataBase mysql配置
type DataBase struct {
	Drive    string `yaml:"drive"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Database string `yaml:"database"`
}

// InitConfig 初始化服务器配置
func InitConfig() *AppConfig {
	var config *AppConfig
	file, err := os.Open("./src/conf.yaml")
	if err != nil {
		panic(err.Error())
	}
	// 解析
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err.Error())
	}
	return config
}
