package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var ProjectConfig *Config

type Config struct {
	Database struct {
		Host   string `yaml:"Host"`
		Port   string `yaml:"Port"`
		User   string `yaml:"Username"`
		Pass   string `yaml:"Password"`
		DBName string `yaml:"DBName"`
	} `yaml:"Database"`
	Volunteer struct {
		TwtKey    string `yaml:"jwt_key"`
		JwtExpiry int    `yaml:"jwt_expiry"`
	} `yaml:"VolunteerConfig"`
	AliyunOSS struct {
		AccessKeyId     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		ObjectKey       string `yaml:"objectKey"`
		BucketName      string `yaml:"bucketName"`
		Endpoint        string `yaml:"endpoint"`
		Area            string `yaml:"area"`
	} `yaml:"AliyunOSSConfig"`
}

func LoadConfig() {
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("打开配置文件出现错误:", err)
	}
	err = yaml.Unmarshal(data, &ProjectConfig)
	if err != nil {
		log.Fatal("解析配置文件失败:", err)
	}
}
