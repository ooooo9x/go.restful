package model

import "time"

type Test struct {
	Id        int64
	Name      string
	CreatTime time.Time `xorm:"created"`
}
