package auth

import (
	"errors"

	"qms.mgmt.api/base/log"

	"github.com/gofrs/uuid"
	"qms.mgmt.api/base/redis"
)

// GetToken 进行用户校验同时生成token，token的生成规则为username + '@@' + uuid
// error为nil时success，否则fail
func GetToken(username string, password string) (string, error) {
	defer log.Logger.Sync()
	log.Logger.Debug("start GetToken")

	// 参数校验
	if username == "" || password == "" {
		log.Logger.Warn("username or password is empty")
		return "", errors.New("username or password is empty")
	}

	// 生成uuid对象
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	uuidstr := u4.String()
	uuidstr = username + "@@" + uuidstr
	log.Logger.Debug("token-->" + uuidstr)

	// 将token存储到redis
	err0 := redis.RedisCache.Set(uuidstr, username, redis.EXPIRES_DEFAULT)
	if err0 != nil {
		log.Logger.Error("error-->" + err0.Error())
		return "", err0
	}
	return uuidstr, nil
}

// ValidToken token校验，error为nil校验成功，否则失败
func ValidToken(token string) error {
	//参数校验
	if token == "" {
		log.Logger.Warn("token is empty")
		return errors.New("token is empty")
	}

	var uuidstr string
	err := redis.RedisCache.Get(token, &uuidstr)
	if err != nil {
		log.Logger.Error("token is fail-->" + err.Error())
		return err
	}

	return nil
}
