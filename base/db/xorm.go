package db

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"qms.mgmt.api/base/log"
)

var Engine *xorm.Engine

func Init() error {
	log.Logger.Info("start init db")
	eg, err := xorm.NewEngine("mysql", "root:root123@tcp(10.221.164.6:3306)/qms?charset=utf8")
	if err != nil {
		return err
	}
	eg.SetMaxIdleConns(5)  //设置连接池的空闲数大小
	eg.SetMaxOpenConns(20) //设置最大打开连接数

	//日志打印SQL
	log.Logger.Info("set db sql --> sql.log")
	eg.ShowSQL(true)
	eg.Logger().SetLevel(core.LOG_DEBUG)
	f, err := os.Create("sql.log")
	if err != nil {
		log.Logger.Info("ping err-->" + err.Error())
		return err
	}
	eg.SetLogger(xorm.NewSimpleLogger(f))

	//连接测试
	if err := eg.Ping(); err != nil {
		log.Logger.Info("ping err-->" + err.Error())
		return err
	}

	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	// eg.SetTableMapper(core.SnakeMapper{})

	Engine = eg
	log.Logger.Info("db init ok!")
	return nil
}
