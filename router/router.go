package router

import (
	"qms.mgmt.api/base/log"

	"qms.mgmt.api/base/restful"
	"qms.mgmt.api/handler/exampleHandler"
	"qms.mgmt.api/handler/tokenHandler"

	"github.com/gin-gonic/gin"
)

// StartServer 启动web服务器，并配置路由
func StartServer() {
	defer log.Logger.Sync()
	log.Logger.Info("start StartServer")

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(AuthFilter())

	log.Logger.Info("start init router")
	// 生成事例group：example
	v1 := router.Group("/example")
	rest.CRUD(v1, "/example0", new(exampleHandler.ExampleHandler))
	rest.CRUD(v1, "/example0", new(exampleHandler.ExampleHandler1))

	// 生成认证类group：auth
	auth := router.Group("/auth")
	// 获取令牌token
	rest.CRUD(auth, "/token", new(tokenHandler.TokenHandler))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// router.Run() //listen and serve on 0.0.0.0:8080
	router.Run(":4000") //for a hard coded port
}
