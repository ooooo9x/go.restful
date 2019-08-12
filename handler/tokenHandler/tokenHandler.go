package tokenHandler

import (
	"net/http"

	"qms.mgmt.api/base/log"

	"github.com/gin-gonic/gin"
	"qms.mgmt.api/base/config"
	"qms.mgmt.api/service/auth"
)

// TokenHandler 获取令牌handler类
type TokenHandler struct{}

// Login 登陆请求input对象
type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateHandle CreateHandle 实现
func (e *TokenHandler) CreateHandle(c *gin.Context) {
	defer log.Logger.Sync()
	log.Logger.Debug("start TokenHandler")

	// 获取input数据并格式校验
	var loginJson Login
	if err := c.ShouldBindJSON(&loginJson); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": config.CODE_VALIDATIONFAIL,
			"mes":  err.Error(),
		})
		return
	}

	// 校验并获取token
	token, err := auth.GetToken(loginJson.Username, loginJson.Password)
	if err != nil {
		log.Logger.Info("TokenHandler fail!")
		c.JSON(http.StatusOK, gin.H{
			"code": config.CODE_AUTHFAIL,
			"mes":  err.Error(),
		})
		return
	}

	// 在header头和cookie里设置token
	c.Header("token", token)
	// c.Header("Set-Cookie", "token="+token+"; Path=/;HttpOnly")
	c.SetCookie("token", token, 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code": config.CODE_SUCCESS,
		"mes":  "TokenHandler success!",
	})

	log.Logger.Debug("TokenHandler success!")
}
