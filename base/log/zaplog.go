package log

import (
	"go.uber.org/zap"
	"qms.mgmt.api/base/config"
)

// Logger 全局logger对象
var Logger *zap.Logger

// InitZapLog 初始化zaplog日志
func InitZapLog() error {
	//生成全局logger
	LoggerVar, err := config.Config.ZapLogConfig.Build()
	if err != nil {
		return err
	}
	Logger = LoggerVar

	defer Logger.Sync()

	Logger.Info("logger construction succeeded")
	return nil
}
