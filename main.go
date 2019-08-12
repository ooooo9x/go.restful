package main

import (
	"log"
	"os"

	"qms.mgmt.api/base/config"
	"qms.mgmt.api/base/db"
	zaplog "qms.mgmt.api/base/log"
	"qms.mgmt.api/base/redis"
	"qms.mgmt.api/router"
)

func main() {
	log.Println("system starting ...")

	//加载配置文件
	err := config.InitConif("")
	if err != nil {
		log.Println("config file load fail! exit(1)-->" + err.Error())
		os.Exit(1)
	}
	//初始化zaplog得logger对象
	log.Println("InitZapLog ...")
	err = zaplog.InitZapLog()
	if err != nil {
		log.Println("log init fail! exit(1)-->" + err.Error())
		os.Exit(1)
	}
	log.Println("Init DB ...")
	err = db.Init()
	if err != nil {
		log.Println("db init fail! exit(1)-->" + err.Error())
		os.Exit(1)
	}
	//初始化redis池
	log.Println("NewRedisCache ...")
	redis.NewRedisCache(config.Config.RedisHost, config.Config.RedisPassword, config.Config.RedisExpires)
	//启动restful服务
	log.Println("StartServer ...")
	router.StartServer()

	log.Println("server exit,bye!")
}
