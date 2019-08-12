package rest

import (
	"reflect"

	"qms.mgmt.api/base/log"

	"github.com/gin-gonic/gin"
)

// CreateEnterPoint POST请求接入点interface
type CreateEnterPoint interface {
	CreateHandle(*gin.Context)
}

// ListEnterPoint GET获取列表请求接入点interface
type ListEnterPoint interface {
	ListHandle(*gin.Context)
}

// GetEnterPoint GET获取单对象请求接入点interface
type GetEnterPoint interface {
	GetHandle(*gin.Context)
}

// UpdateEnterPoint PUT请求接入点interface
type UpdateEnterPoint interface {
	UpdateHandle(*gin.Context)
}

// DeleteEnterPoint DELETE请求接入点interface
type DeleteEnterPoint interface {
	DeleteHandle(*gin.Context)
}

// CRUD restful路由配置方法
func CRUD(group *gin.RouterGroup, path string, resource interface{}) {
	defer log.Logger.Sync()
	log.Logger.Info("CRUD start ...type-->" + reflect.TypeOf(resource).Name())

	if resource, ok := resource.(CreateEnterPoint); ok {
		log.Logger.Info("CRUD-->CreateEnterPoint   path-->" + path)
		group.POST(path, resource.CreateHandle)
	}
	if resource, ok := resource.(ListEnterPoint); ok {
		log.Logger.Info("CRUD-->ListEnterPoint   path-->" + path)
		group.GET(path, resource.ListHandle)
	}
	if resource, ok := resource.(GetEnterPoint); ok {
		log.Logger.Info("CRUD-->GetEnterPoint   path-->" + path)
		group.GET(path+"/:id", resource.GetHandle)
	}
	if resource, ok := resource.(UpdateEnterPoint); ok {
		log.Logger.Info("CRUD-->UpdateEnterPoint   path-->" + path)
		group.PUT(path+"/:id", resource.UpdateHandle)
	}
	if resource, ok := resource.(DeleteEnterPoint); ok {
		log.Logger.Info("CRUD-->DeleteEnterPoint   path-->" + path)
		group.DELETE(path+"/:id", resource.DeleteHandle)
	}
}
