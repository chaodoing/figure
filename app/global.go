package app

import (
	"encoding/xml"
	`errors`
	"fmt"
	"io"
	`log`
	"os"
	"path"
	`strings`
	"time"
	
	"github.com/go-redis/redis"
	"github.com/gookit/goutil/fsutil"
	"github.com/lestrrat-go/strftime"
	encoder "github.com/zwgblue/yaml-encoder"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	`gorm.io/gorm/logger`
)

var dbLogLevel = map[string]logger.LogLevel{
	`silent`: logger.LogLevel(1),
	`error`:  logger.LogLevel(2),
	`warn`:   logger.LogLevel(3),
	`info`:   logger.LogLevel(4),
}

type (
	// Logger 日志配置
	Logger struct {
		Console bool   `json:"console" xml:"console" yaml:"Console" comment:"是否输出到控制台"` // Console 是否输出到控制台
		File    string `json:"file" xml:"file" yaml:"File" comment:"日志文件"`                  // File 日志文件
		Level   string `json:"level" xml:"level" yaml:"Level" comment:"日志等级"`               // Level 日志等级
	}
	// Resource 静态资源文件配置
	Resource struct {
		Url string `json:"url" xml:"url" yaml:"Url" comment:"访问路径"`     // Url 访问路径
		Dir string `json:"dir" xml:"dir" yaml:"Dir" comment:"静态资源位置"` // Dir 静态资源位置
	}
	// Template 模板配置
	Template struct {
		Delimit []string `json:"delimit" xml:"delimit" delim:"|" yaml:"Delimit" comment:"模板使用分隔符"` // Delimit 模板使用分隔符
		Dir     string   `json:"dir" xml:"dir" yaml:"Dir" comment:"模板目录位置"`                         // Dir 模板目录位置
		Ext     string   `json:"ext" xml:"ext" yaml:"Ext" comment:"模板文件扩展名称"`                     // Ext 模板文件扩展名称
	}
	// Service iris应用配置
	Service struct {
		Favicon     string     `json:"favicon" xml:"favicon" yaml:"Favicon" comment:"网站图标配置"`                             // Favicon 网站图标配置
		Port        uint16     `json:"port" xml:"port" yaml:"Port" comment:"监听端口"`                                          // Port 监听端口
		Host        string     `json:"host" xml:"host" yaml:"Host" comment:"监听主机"`                                          // Host 监听主机
		CrossDomain bool       `json:"cross_domain" xml:"crossDomain" yaml:"CrossDomain" comment:"允许跨域"`                    // CrossDomain 允许跨域
		Log         *Logger    `json:"log" xml:"log" yaml:"Log" comment:"日志配置 level:[disable fatal error warn info debug]"` // Log 日志配置
		Template    *Template  `json:"template" xml:"template" yaml:"Template" comment:"模板目录配置"`                          // Template 模板目录配置
		Resources   []Resource `json:"resources" xml:"resources" yaml:"Resources" comment:"允许跨域"`                           // Resources 静态资源文件配置
	}
	// Redis redis配置
	Redis struct {
		Host string `json:"host" xml:"host" yaml:"Host" comment:"连接主机"` // Host 连接主机
		Port uint16 `json:"port" xml:"port" yaml:"Port" comment:"连接端口"` // Port 连接端口
		Db   int    `json:"db" xml:"db" yaml:"Db" comment:"数据库索引"`     // Db 数据库索引
		Auth string `json:"auth" xml:"auth" yaml:"Auth" comment:"连接密码"` // Auth 连接密码
		TTL  uint64 `json:"ttl" xml:"ttl" yaml:"TTL" comment:"缓存时长"`    // TTL 缓存时长
	}
	// MySQL mysql配置
	MySQL struct {
		Host     string  `json:"host" xml:"host" yaml:"Host" comment:"连接主机"`                                      // Host 连接主机
		Port     uint16  `json:"port" xml:"port" yaml:"Port" comment:"连接端口"`                                      // Port 连接端口
		Name     string  `json:"name" xml:"name" yaml:"Name" comment:"数据库名称"`                                    // Name 数据库名称
		Username string  `json:"username" xml:"username" yaml:"Username" comment:"连接用户名"`                        // Username 连接用户名
		Password string  `json:"password" xml:"password" yaml:"Password" comment:"连接密码"`                          // Password 连接密码
		Charset  string  `json:"charset" xml:"charset" yaml:"Charset" comment:"连接字符集"`                           // Charset 连接字符集
		Logger   *Logger `json:"logger" xml:"logger" yaml:"Logger" comment:"日志配置 level:[silent error warn info]"` // Logger 日志配置
	}
)

// Global 结构体包含了应用全局配置，包括服务配置、Redis配置和MySQL数据库配置。
type Global struct {
	XMLName xml.Name      `json:"-" xml:"root" yaml:"-"`
	Service *Service      `json:"service" xml:"service" yaml:"Service" comment:"服务配置"` // Service 结构体用于定义服务配置。
	Redis   *Redis        `json:"redis" xml:"redis" yaml:"Redis" comment:"服务配置"`       // Redis 结构体用于定义Redis服务配置。
	MySQL   *MySQL        `json:"mysql" xml:"mysql" yaml:"MySQL" comment:"数据库配置"`     // MySQL 结构体用于定义MySQL数据库配置。
	db      *gorm.DB      `json:"-" xml:"-" yaml:"-"`                                      // db 是一个 *gorm.DB 类型，用于存储GORM数据库会话实例，不暴露给JSON、XML或YAML序列化。
	rdx     *redis.Client `json:"-" xml:"-" yaml:"-"`                                      // rdx 是一个 *redis.Client 类型，用于存储Redis客户端实例，不暴露给JSON、XML或YAML序列化。
}

func (g *Global) Rdx() *redis.Client {
	return g.rdx
}

// GlobalDefault 返回一个预设的全局配置对象。
//
// 返回值 Global 包括了服务配置（Service）、MySQL数据库配置（MySQL）、Redis缓存配置（Redis）等。
func GlobalDefault() Global {
	return Global{
		Service: &Service{
			// 服务配置包括网站图标路径、监听端口、主机地址、跨域设置及日志配置。
			Favicon:     "${DIR}/resources/favicon.ico",
			Port:        8080,
			Host:        "127.0.0.1",
			CrossDomain: true,
			Log: &Logger{
				// 日志配置包括控制台输出、文件输出路径及日志级别。
				Console: true,
				File:    "${DIR}/logs/app-%Y-%m-%d.log",
				Level:   "debug",
			},
			// 资源配置包括静态文件的URL和目录路径。
			Resources: []Resource{
				{Url: "/upload", Dir: "${DIR}/resources/upload"},
				{Url: "/static", Dir: "${DIR}/resources/static"},
			},
			// 模板配置包括模板引擎的标识符、目录路径和文件扩展名。
			Template: &Template{
				Delimit: []string{"{{", "}}"},
				Dir:     "${DIR}/resources/template",
				Ext:     ".html",
			},
		},
		MySQL: &MySQL{
			// MySQL数据库配置包括主机地址、端口、数据库名、用户名、密码及日志配置。
			Host:     "127.0.0.1",
			Port:     3306,
			Name:     "dbName",
			Username: "root",
			Password: "password",
			Charset:  "utf8mb4",
			Logger: &Logger{
				// 数据库操作的日志配置，包括控制台输出、文件输出路径及日志级别。
				Console: true,
				File:    "${DIR}/logs/mysql-%Y-%m-%d.log",
				Level:   "info",
			},
		},
		Redis: &Redis{
			// Redis缓存配置包括主机地址、端口、数据库索引、认证密码及TTL（过期时间）。
			Host: "127.0.0.1",
			Port: 6379,
			Db:   0,
			Auth: "password",
			TTL:  648000,
		},
	}
}

// Dialect 方法用于构造并返回MySQL数据库的连接字符串。
//
// 返回值:
// schema: 格式化后的MySQL连接字符串，包含了用户名、密码、主机、端口、数据库名和字符集等信息。
func (c *MySQL) Dialect() (schema string) {
	// 构造MySQL连接字符串
	schema = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Name, c.Charset)
	return
}

// Log 为指定名称创建一个日志文件。
//
// 参数:
//
//	name string - 日志文件的名称。
//
// 返回值:
//
//	*lumberjack.Logger - 指向创建的日志文件的指针。
//	error - 如果创建过程中遇到错误，则返回错误信息；否则为nil。
func (g *Global) Log(name string) (drive *lumberjack.Logger, err error) {
	// 获取日志文件所在目录
	dir := path.Dir(name)
	if !fsutil.PathExists(dir) {
		if err = fsutil.Mkdir(dir, 0755); err != nil {
			return
		}
	}
	var p *strftime.Strftime
	// 使用指定的文件名创建一个格式化器
	p, err = strftime.New(name)
	// 配置日志文件参数
	drive = &lumberjack.Logger{
		Filename:   p.FormatString(time.Now()), // 设置日志文件名，使用时间戳区分
		MaxSize:    4,                          // 最大尺寸，单位MB
		MaxAge:     1,                          // 最长保留时间，单位天
		MaxBackups: 31,                         // 最多备份文件数量
		LocalTime:  true,                       // 使用本地时间
		Compress:   true,                       // 是否压缩备份文件
	}
	return
}

// MakeEnv 生成环境配置文件
// 此函数创建一个环境配置文件（.env.yml），基于全局配置实例g的内容。
// 它首先使用一个带注释的编码器将g编码为YAML格式，然后将编码后的值写入到环境配置文件中。
//
// 参数:
// g *Global: 全局配置实例，包含了需要被编码为环境配置文件的内容。
//
// 返回值:
// error: 如果在编码或写入文件过程中遇到错误，将返回一个error；否则返回nil。
func (g *Global) MakeEnv() (err error) {
	// 使用带注释的编码器对全局配置g进行编码
	value, err := encoder.NewEncoder(g, encoder.WithComments(encoder.CommentsOnHead)).Encode()
	if err != nil {
		return
	}
	// 将编码后的配置写入到环境配置文件中
	_, err = fsutil.PutContents(os.ExpandEnv("${DIR}/.env.yml"), value)
	return
}

// parseDir 方法解析配置中的环境变量
// 对 Global 结构体中 Service 字段的 Favicon, Log.File, Resources[].Dir, 和 Template.Dir
// 属性进行环境变量展开，以便在不同环境下使用不同的路径。
// 参数:
//
//	无
//
// 返回值:
//
//	*Global: 返回经过环境变量解析后的 Global 实例指针。
func (g *Global) parseDir() *Global {
	// 展开 Service 字段中 Favicon 和 Log.File 的环境变量
	g.Service.Favicon = os.ExpandEnv(g.Service.Favicon)
	g.Service.Log.File = os.ExpandEnv(g.Service.Log.File)
	// 遍历 Resources 列表，展开每个资源的 Dir 字段中的环境变量
	for i, resource := range g.Service.Resources {
		resource.Dir = os.ExpandEnv(resource.Dir)
		g.Service.Resources[i] = resource
	}
	// 展开 Template.Dir 字段中的环境变量
	g.Service.Template.Dir = os.ExpandEnv(g.Service.Template.Dir)
	return g
}

// LoadEnv 加载环境配置文件
// 该函数负责从指定路径加载.env.yml文件，并将其内容解析到Global结构体中。
// 如果解析成功，将返回一个指向Global结构体的指针，否则返回nil。
// 参数：
//
//	g *Global: 指向Global结构体的指针，用于存储解析后的环境配置。
//
// 返回值：
//
//	*Global: 解析成功时返回指向Global结构体的指针，失败时返回nil。
func (g *Global) LoadEnv() *Global {
	// 从环境变量扩展的路径中读取.env.yml文件的内容
	value := fsutil.GetContents(os.ExpandEnv("${DIR}/.env.yml"))
	// 将读取到的内容解析到g中
	err := yaml.Unmarshal(value, g)
	if err != nil {
		// 解析失败时返回nil
		return nil
	}
	// 解析成功时，进一步处理并返回
	return g.parseDir()
}

// Db 数据库初始化
func (g *Global) Db() (db *gorm.DB, err error) {
	if g.db != nil {
		return g.db, nil
	}
	var write io.Writer
	write, err = g.Log(g.MySQL.Logger.File)
	if err != nil {
		return
	}
	if g.MySQL.Logger.Console {
		write = io.MultiWriter(write, os.Stdout)
	}
	// 创建一个新的logger实例。
	//
	// write: 提供一个写入接口，用于logger写入日志。
	// g.MySQL.Logger.Level: 从外部配置获取的日志级别。
	// 返回值: 初始化后的logger实例。
	logs := logger.New(log.New(write, "", log.LstdFlags|log.Ldate|log.Ltime), logger.Config{
		SlowThreshold:             0,                                // 慢查询阈值设置为0
		Colorful:                  false,                            // 不使用彩色日志
		IgnoreRecordNotFoundError: false,                            // 不忽略记录未找到的错误
		ParameterizedQueries:      false,                            // 不使用参数化查询日志
		LogLevel:                  dbLogLevel[g.MySQL.Logger.Level], // 根据配置设定日志级别
	})
	g.db, err = gorm.Open(mysql.Open(g.MySQL.Dialect()), &gorm.Config{
		SkipDefaultTransaction: true,  // SkipDefaultTransaction 跳过默认事务
		FullSaveAssociations:   true,  // FullSaveAssociations 在创建或更新时，是否更新关联数据
		Logger:                 logs,  // Logger 日志接口，用于实现自定义日志
		DryRun:                 false, // DryRun 生成 SQL 但不执行，可以用于准备或测试生成的 SQL
		PrepareStmt:            true,  // PrepareStmt 是否禁止创建 prepared statement 并将其缓存
		AllowGlobalUpdate:      false, // AllowGlobalUpdate 是否允许全局 update/delete
		QueryFields:            true,  // QueryFields 执行查询时，是否带上所有字段
	})
	return g.db, err
}

// Rds redis初始化
func (g *Global) Rds() (rdx *redis.Client, err error) {
	if g.rdx != nil {
		return g.rdx, nil
	}
	g.rdx = redis.NewClient(&redis.Options{
		DB:       g.Redis.Db,                                       // DB 数据库索引
		Addr:     fmt.Sprintf("%s:%d", g.Redis.Host, g.Redis.Port), // Addr 连接地址
		Password: g.Redis.Auth,                                     // Password 连接密码
	})
	pong, err := g.rdx.Ping().Result()
	if err != nil {
		return
	}
	if err != nil || !strings.EqualFold(pong, "PONG") {
		err = errors.New("redis connection error")
		return
	}
	return g.rdx, err
}
