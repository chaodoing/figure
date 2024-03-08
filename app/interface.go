package app

import (
	`github.com/go-redis/redis`
	`gopkg.in/natefinch/lumberjack.v2`
	`gorm.io/gorm`
)

type GlobalInterface interface {
	Db() (db *gorm.DB, err error)
	Rds() (rdx *redis.Client, err error)
	Log() (logger *lumberjack.Logger, err error)
	Auth() (auth interface{}, err error)
}
