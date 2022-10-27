package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 项目全局配置

var Conf = &config{}

func InitConfig(path string) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Init project config error, %s", err.Error())
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatalf("Parse project config error, %s", err.Error())
	}
}

type config struct {
	AppName     string `mapstructure:"app_name"`
	LogLevel    string `mapstructure:"log_level"`     // 日志级别
	Port        int    `mapstructure:"port"`          // 项目运行端口
	GrpcPort    int    `mapstructure:"grpc_port"`     // grpc运行端口
	MaxFileSize int64  `mapstructure:"max_file_size"` // 最大文件上传大小，byte
	Etcd        etcd   `mapstructure:"etcd"`
}

type etcd struct {
	Addr      string `mapstructure:"addr"`
	HeartBeat int64  `mapstructure:"heart_beat"`
}
