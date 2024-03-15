package app

import (
	`io`
	`time`
	
	`github.com/go-redis/redis`
	`gorm.io/gorm`
)

// EventInterface 定义了一个事件处理的接口。
type EventInterface interface {
	// DbInit 初始化数据库连接。
	DbInit()
	
	// RedisInit 初始化Redis连接。
	RedisInit()
	
	// LogInit 初始化日志系统。
	// name: 日志文件的名称。
	LogInit(name string)
	
	// EnvInit 初始化环境配置。
	EnvInit(global Global)
}

// GlobalInterface 定义了全局接口，包含了数据库、Redis客户端、日志记录器、授权信息的获取方法，以及环境变量的制作和加载方法。
type GlobalInterface interface {
	// Db 返回一个初始化好的Gorm数据库实例和可能发生的错误。
	Db() (db *gorm.DB, err error)
	// Rds 返回一个初始化好的Redis客户端实例和可能发生的错误。
	Rds() (rdx *redis.Client, err error)
	// Log 根据提供的名称和控制台标志，创建并返回一个日志记录器实例和可能发生的错误。
	Log(name string, console bool) (write io.Writer, err error)
	// Auth 返回一个初始化好的Authorization实例和可能发生的错误。
	Auth() (auth Authorization, err error)
	// MakeEnv 用于创建和配置环境变量，返回可能发生的错误。
	MakeEnv() (err error)
	// LoadEnv 用于加载环境变量，并返回一个Global实例，包含加载的全局配置信息。
	LoadEnv() (global Global)
}

// CacheInterface 接口定义了缓存的基本操作。
type CacheInterface interface {
	// Set 设置缓存项。
	// key: 缓存项的键。
	// value: 缓存项的值。
	// expiration: 缓存项的过期时间。
	// 返回值: 设置成功返回 nil，失败返回 error。
	Set(key string, value string, expiration time.Duration) (err error)
	
	// Get 从缓存中获取项的值。
	// key: 缓存项的键。
	// 返回值: 缓存项的值，如果不存在则返回空字符串。
	Get(key string) string
	
	// Has 检查缓存中是否存在指定键的缓存项。
	// key: 缓存项的键。
	// 返回值: 存在返回 true，不存在返回 false。
	Has(key string) bool
	
	// Del 删除缓存项。
	// key: 缓存项的键。
	// 返回值: 删除成功返回 nil，失败返回 error。
	Del(key string) (err error)
	
	// Clear 清除缓存中的所有项。
	// 返回值: 清除成功返回 nil，失败返回 error。
	Clear() error
	
	// NewKey 生成新的缓存键，同时转移旧键的值到新键，并设置新键的过期时间。
	// oldKey: 原始缓存项的键。
	// newKey: 新的缓存项键。
	// ttl: 新缓存项的过期时间。
	// 返回值: 生成成功返回 nil，失败返回 error。
	NewKey(oldKey, newKey string, ttl time.Duration) (err error)
}
