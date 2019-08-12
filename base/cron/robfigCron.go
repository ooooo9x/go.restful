package cron

import (
	"log"

	"github.com/robfig/cron"
)

// CreateCronByDay robfigCron实例
func CreateCronByDay() {
	c := cron.New()
	specTime := "20 0 0 * * ?"
	c.AddFunc(specTime, func() {
		log.Println("任务开始")
	})
	c.Start()
}
