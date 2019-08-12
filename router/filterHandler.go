package router

import (
	"errors"
	"net/http"

	"qms.mgmt.api/base/log"

	"github.com/gin-gonic/gin"
	"qms.mgmt.api/base/config"
	"qms.mgmt.api/service/auth"
)

// AuthFilter 获取权限校验过滤器
func AuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer log.Logger.Sync()
		err := filterHandle(c)
		if err != nil {
			log.Logger.Warn("auth fail -->" + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": config.CODE_AUTHFAIL,
				"mes":  err.Error(),
			})
			return
		}
	}
}

// FilterHandle 接口实现
func filterHandle(c *gin.Context) error {
	defer log.Logger.Sync()
	log.Logger.Debug("start FilterHandle")

	path := c.Request.URL.Path
	log.Logger.Debug("path-->" + path)

	//过滤免校验的path
	if _, ok := config.Config.AuthExclusion[path]; ok {
		return nil
	}

	//优先从cookie读取token值，并验证
	token, err := c.Cookie("token")
	log.Logger.Debug("Cookie token-->" + token)
	if err == nil {
		err := auth.ValidToken(token)
		if err == nil {
			log.Logger.Debug("Cookie token is ok!")
			return nil
		}
	} else {
		log.Logger.Info("Cookie err-->" + err.Error())
	}

	//cookie验证不通过时，从header读取token值，并验证
	token = c.GetHeader("token")
	log.Logger.Debug("Header token-->" + token)
	if token != "" {
		err := auth.ValidToken(token)
		if err == nil {
			log.Logger.Debug("Header token is ok!")
			return nil
		}
		log.Logger.Warn("ValidToken error-->" + err.Error())
	} else {
		log.Logger.Info("Header err-->" + err.Error())
	}

	return errors.New("auth fail!")
}
