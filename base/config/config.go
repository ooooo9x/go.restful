package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	// CODE_SUCCESS 操作成功
	CODE_SUCCESS = 0
	// CODE_FAIL 操作失败
	CODE_FAIL = 1
	// CODE_AUTHFAIL 鉴权失败
	CODE_AUTHFAIL = 2
	// CODE_VALIDATIONFAIL 输入解析失败
	CODE_VALIDATIONFAIL = 3
)

// Configuration 配置对象
type Configuration struct {
	RedisHost     string          `json:"redis.host"`
	RedisPassword string          `json:"redis.password"`
	RedisExpires  time.Duration   `json:"redis.expires"`
	AuthToken     bool            `json:"auth.token"`
	AuthExclusion map[string]bool `json:"auth.exclusion"`
	ZapLogConfig  zap.Config      `json:"zaplog.config"`
}

// Config 定义一个Configuration的全局变量
var Config Configuration

// InitConif 加载conf.json配置文件到Configuration对象
func InitConif(configFile string) error {
	log.Println("configFile -->" + configFile)
	if configFile == "" {
		configFile = "conf.json"
	}

	file, _ := os.Open(configFile)
	defer file.Close()

	//初始化config
	Config = Configuration{RedisHost: "127.0.0.1:6379",
		RedisPassword: "admin",
		RedisExpires:  time.Duration(86400000000000),
		AuthToken:     true,
		AuthExclusion: map[string]bool{},
	}
	Config.ZapLogConfig = zap.NewDevelopmentConfig()

	//解析json文件
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	//检查对象并输出
	log.Println("check config value: ")
	log.Println(Config)

	return nil
}
