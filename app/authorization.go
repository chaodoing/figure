package app

import (
	`encoding/json`
	`errors`
	`strings`
	`time`
	
	`github.com/chaodoing/figure/encrypt`
	`github.com/go-redis/redis`
	`github.com/gookit/goutil/arrutil`
	`github.com/kataras/iris/v12`
)

const (
	Basic  = "Basic "
	Bearer = "Bearer "
)

type (
	// refresh 结构体用于表示刷新认证信息
	refresh struct {
		Tokenization string    // Tokenization 访问令牌
		Expire       time.Time // Expire 新认证字符串有效期
	}
	
	// Authorization 结构体用于管理认证和刷新令牌的逻辑
	Authorization struct {
		cache        *Cache       // cache 用于存储认证信息的缓存
		ctx          iris.Context // ctx 表示当前的上下文环境
		isRefresh    bool         // isRefresh 标记是否是刷新令牌操作
		Tokenization string       // Tokenization 访问令牌
		TTL          uint64       // TTL 存储有效期，单位通常为秒
		Refresh      refresh      // Refresh 包含新的认证字符串和其有效期
	}
)

// NewAuthorization 创建一个新的Authorization实例。
//
// 参数:
//   rds *redis.Client: 一个指向Redis客户端的指针，用于缓存管理。
//   TTL uint64: 缓存条目的时间-to-live(TTL)，单位为秒。
//
// 返回值:
//   Authorization: 返回一个初始化好的Authorization实例。
//
func NewAuthorization(rds *redis.Client, TTL uint64) Authorization {
	return Authorization{
		cache: &Cache{rdx: rds}, // 初始化cache字段，使用传入的redis.Client.
		TTL:   TTL,              // 设置TTL字段为传入的值.
	}
}

// refresh 方法用于刷新授权信息。
// 返回值包含刷新后的令牌和过期时间。
func (a Authorization) refresh() refresh {
	// 生成新的令牌并设置过期时间
	return refresh{
		Tokenization: encrypt.UUID(),                                     // 使用加密库生成一个新的UUID作为令牌
		Expire:       time.Now().Add(time.Duration(a.TTL) * time.Second), // 设置令牌的过期时间为当前时间加上TTL指定的时间长度
	}
}

// init 初始化函数
// 该函数用于初始化Authorization对象，通过解析上下文中的刷新令牌来决定是否需要刷新令牌。
// 参数:
// - ctx iris.Context: 传递当前的iris上下文环境，用于获取请求头等信息。
// 返回值:
// - Authorization: 返回初始化后的Authorization对象。
func (a Authorization) init(ctx iris.Context) Authorization {
	// 从请求头中获取刷新令牌字符串
	refreshString := ctx.GetHeader("Refresh-Token")
	a.ctx = ctx                // 保存当前的上下文环境至Authorization对象中
	a.Tokenization = a.token() // 初始化令牌
	// 判断是否需要进行刷新令牌的操作
	a.isRefresh = arrutil.NotIn(refreshString, []string{"false", "0", "off"})
	if a.isRefresh {
		// 如果需要刷新令牌，则执行刷新操作
		a.Refresh = a.refresh()
	}
	a.header() // 处理相关头部信息
	return a
}

// header 方法设置跨域资源共享（CORS）相关的响应头，并根据是否是刷新令牌来添加相应的令牌信息。
// 该方法不接受参数，也不返回值。
// 其中，a 是 Authorization 类型的实例，该实例包含处理请求所需的上下文信息和认证相关的数据。
func (a Authorization) header() {
	// 设置允许的请求头，这些头部信息会在浏览器中暴露给前端。
	a.ctx.Header("Access-Control-Allow-Headers", "Refresh-Token, Accept-Version, Authorization, Access-Token, Language, Access-Control-Allow-Methods, Access-Control-Allow-Origin, Cache-Control, Content-Type, if-match, if-modified-since, if-none-match, if-unmodified-since, X-Requested-With")
	// 设置允许浏览器访问的响应头，这些头部信息对于前端通过 AJAX 请求获取非常重要。
	a.ctx.Header("Access-Control-Expose-Headers", "Authorization, Access-Token, Refresh-Token, Refresh-Expires")
	// 根据是否是刷新令牌，设置相应的令牌信息到响应头。
	if a.isRefresh {
		// 如果是刷新令牌，则设置刷新令牌和过期时间。
		a.ctx.Header("Refresh-Token", a.Refresh.Tokenization)
		a.ctx.Header("Refresh-Expires", a.Refresh.Expire.Format("2006-01-02 15:04:05"))
	} else {
		// 如果不是刷新令牌，则设置普通令牌信息。
		a.ctx.Header("Refresh-Token", a.Tokenization)
	}
}

// token 方法用于获取授权令牌。
// 如果上下文（ctx）不为空，则尝试从 Authorization 和 Accept-Token 头部信息中提取令牌，
// 如果这些头部信息中没有找到有效的令牌，则生成一个新的UUID作为令牌返回。
// 如果上下文（ctx）为空，直接生成一个新的UUID作为令牌返回。
// 参数:
//   - a: Authorization 类型的实例
// 返回值:
//   - string: 令牌的值
func (a Authorization) token() (value string) {
	if a.ctx != nil {
		// 尝试从 Authorization 头部信息中提取令牌，移除前缀 "Bearer" 或 "Basic"。
		value = strings.TrimPrefix(strings.TrimPrefix(a.ctx.GetHeader("Authorization"), Bearer), Basic)
		// 如果提取到的值不为空，则返回该值。
		if len(value) != 0 {
			return
		}
		// 尝试从 Accept-Token 头部信息中提取令牌，移除前缀 "Bearer" 或 "Basic"。
		value = strings.TrimPrefix(strings.TrimPrefix(a.ctx.GetHeader("Accept-Token"), Bearer), Basic)
		// 如果提取到的值不为空，则返回该值。
		if len(value) != 0 {
			return
		}
		// 如果头部信息中没有找到有效的令牌，生成一个新的UUID作为令牌返回。
		return encrypt.UUID()
	}
	// 如果上下文为空，直接生成一个新的UUID作为令牌返回。
	return encrypt.UUID()
}

// SET 方法用于将数据存储到缓存中。
// ctx: iris上下文，用于初始化授权。
// data: 需要存储的数据，任意类型。
// 返回值: 错误信息，如果操作成功，则返回nil。
func (a Authorization) SET(ctx iris.Context, data interface{}) (err error) {
	a.init(ctx) // 使用iris上下文初始化授权对象。
	var value []byte
	value, err = json.Marshal(data) // 将data数据序列化为JSON格式。
	if err != nil {
		return // 如果序列化过程中出现错误，直接返回错误。
	}
	// 将序列化后的数据存储到缓存中，设置过期时间为TTL秒。
	err = a.cache.Set(a.Tokenization, string(value), time.Duration(a.TTL)*time.Second)
	return
}

// GET 用于通过授权进行数据获取。
// ctx: iris上下文，用于处理HTTP请求和响应。
// data: 用于存储从缓存或API获取的数据的接口类型变量。
// 返回值 err: 错误信息，如果操作成功，则为nil。
func (a Authorization) GET(ctx iris.Context, data interface{}) (err error) {
	// 初始化授权信息
	a.init(ctx)
	// 检查缓存中是否存在授权令牌
	if a.cache.Has(a.Tokenization) {
		value := a.cache.Get(a.Tokenization)
		// 如果缓存值不为空，则尝试从缓存中解析数据
		if !strings.EqualFold(value, "") {
			err = json.Unmarshal([]byte(value), &data) // 解析缓存中的数据
			return
		} else {
			// 缓存中数据为空时返回错误
			return errors.New("数据不存在")
		}
	} else {
		// 缓存中无授权令牌时返回错误
		return errors.New("数据不存在")
	}
}

// DEL 从缓存中删除指定的令牌
// @param ctx iris.Context: 上下文信息，用于处理当前请求
// @return err error: 操作过程中可能出现的错误
func (a Authorization) DEL(ctx iris.Context) (err error) {
	a.init(ctx) // 初始化授权信息
	// 检查缓存中是否存在指定的令牌
	if a.cache.Has(a.Tokenization) {
		// 如果存在，则从缓存中删除该令牌
		err = a.cache.Del(a.Tokenization)
	}
	return
}

// MOD 方法用于更新认证的Token。
// 参数:
// - ctx: iris上下文，用于获取请求相关信息和进行响应。
// 返回值:
// - err: 错误信息，如果操作成功，则为nil。
func (a Authorization) MOD(ctx iris.Context) (err error) {
	a.init(ctx) // 初始化认证信息
	// 检查缓存中是否存在Token
	if a.cache.Has(a.Tokenization) {
		// 如果存在，则更新Token的缓存时间
		err = a.cache.NewKey(a.Tokenization, a.Refresh.Tokenization, time.Duration(a.TTL)*time.Second)
	}
	return
}
