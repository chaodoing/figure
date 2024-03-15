package app

import (
	`time`
	
	`github.com/go-redis/redis`
)

type Cache struct {
	rdx *redis.Client
}

// Get 从缓存中获取指定键对应的值。
//
// 参数:
//   key string - 要获取值的键。
//
// 返回值:
//   string - 如果找到对应的键，则返回其值；否则返回空字符串。
func (c *Cache) Get(key string) string {
	// 尝试从缓存中获取指定键的值
	value, err := c.rdx.Get(key).Result() // 使用Redis的Get方法获取值
	if err != nil {
		// 如果获取过程中出现错误，则返回空字符串
		return ""
	}
	// 如果成功获取到值，则返回该值
	return value
}

// Set 在缓存中设置一个键值对，并指定其过期时间。
//
// 参数:
//   key - 要设置的键。
//   value - 键对应的值。
//   expiration - 键的过期时间。
//
// 返回值:
//   err - 设置过程中遇到的错误，如果设置成功则为 nil。
func (c *Cache) Set(key string, value string, expiration time.Duration) (err error) {
	_, err = c.rdx.Set(key, value, expiration).Result() // 尝试在缓存系统中设置键值对并指定过期时间
	return
}

// Has 检查缓存中是否存在指定的键。
//
// 参数:
//   key string - 要检查的键。
//
// 返回值:
//   bool - 如果缓存中存在该键，则返回 true；否则返回 false。
func (c *Cache) Has(key string) bool {
	// 尝试检查指定键是否存在
	v, err := c.rdx.Exists(key).Result() // 调用Redis的Exists方法检查键是否存在
	if err != nil {
		return false // 如果存在错误，认为键不存在
	}
	if v == 0 {
		return false // 如果返回值为0，表示键不存在
	}
	return true // 键存在
}

// Del 从缓存中删除指定的键。
//
// 参数:
//   key string - 需要删除的键。
//
// 返回值:
//   err error - 操作过程中发生的任何错误。
func (c *Cache) Del(key string) (err error) {
	// 尝试从缓存中删除指定的键，并处理操作结果。
	_, err = c.rdx.Del(key).Result() // rdx 是 Redis 客户端的实例
	return
}

// NewKey 用一个新的键替换旧的键，并设置新的过期时间。
// oldKey: 需要被替换的旧键。
// newKey: 用来替换旧键的新键。
// ttl: 新键的过期时间。
// 返回值 err: 操作过程中发生的任何错误。
func (c *Cache) NewKey(oldKey, newKey string, ttl time.Duration) (err error) {
	// 尝试从缓存中获取旧键对应的值
	value, err := c.rdx.Get(oldKey).Result()
	if err != nil {
		return
	}
	// 删除旧键
	err = c.rdx.Del(oldKey).Err()
	if err != nil {
		return
	}
	// 使用新键和获取到的值在缓存中设置一个新的条目，并设置其过期时间
	err = c.rdx.Set(newKey, value, ttl).Err()
	return
}

// Clear 清除缓存
// 此方法用于清空缓存数据库中的所有数据。
// 返回值表示操作是否成功，成功则返回nil，失败则返回错误信息。
func (c *Cache) Clear() error {
	// 调用rdx的FlushDB方法清空数据库，并返回操作结果的错误信息
	return c.rdx.FlushDB().Err()
}
