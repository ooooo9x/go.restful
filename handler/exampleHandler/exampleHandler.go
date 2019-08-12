package exampleHandler

import (
	"net/http"

	"qms.mgmt.api/base/db"
	"qms.mgmt.api/base/log"
	"qms.mgmt.api/model"

	"github.com/gin-gonic/gin"
	"qms.mgmt.api/base/config"
)

// ExampleHandler Example结构体
type ExampleHandler struct{}

// GetHandle 实现
func (e *ExampleHandler) GetHandle(c *gin.Context) {
	defer log.Logger.Sync()
	log.Logger.Debug("start ExampleHandler")

	id := c.Param("id")
	// c.String(http.StatusOK, "Hello %s", id)
	c.JSON(http.StatusOK, gin.H{
		"code": config.CODE_SUCCESS,
		"mes":  "hello " + id,
	})
}

// ExampleHandler Example结构体
type ExampleHandler1 struct{}

// Test input对象
type HttpTest struct {
	Name string `json:"name" binding:"required"`
}

// CreateHandle 实现
func (e *ExampleHandler1) CreateHandle(c *gin.Context) {
	defer log.Logger.Sync()
	log.Logger.Debug("add test object")

	var testJson HttpTest
	if err := c.ShouldBindJSON(&testJson); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": config.CODE_VALIDATIONFAIL,
			"mes":  err.Error(),
		})
		return
	}

	testObj := model.Test{Name: testJson.Name}
	affected, err := db.Engine.Insert(&testObj)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": config.CODE_FAIL,
			"mes":  err.Error(),
		})
		return
	}

	// c.String(http.StatusOK, "Hello %s", id)
	c.JSON(http.StatusOK, gin.H{
		"code": config.CODE_SUCCESS,
		"mes":  "hello " + string(affected),
	})
}
