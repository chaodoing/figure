package app

import (
	`fmt`
	`os`
	`strings`
	
	`github.com/gookit/goutil/fsutil`
	`github.com/kataras/iris/v12`
	`github.com/kataras/iris/v12/hero`
	`github.com/kataras/iris/v12/mvc`
)

type Bootstrap struct {
	Global Global
	m      *mvc.Application
	app    *iris.Application
}

// New 创建一个新的Bootstrap实例。
//
// 参数:
// - file: 指定配置文件的路径。
// - isJson: 标识配置文件的格式是否为JSON。如果是true，则读取JSON格式的配置文件；否则，读取XML格式的配置文件。
// - event: 事件接口，用于在处理配置文件时触发事件。
//
// 返回值:
// - Bootstrap: 返回一个初始化好的Bootstrap实例。
// - error: 如果在创建过程中遇到任何错误，则返回非nil的error。
func New(file string, isJson bool, event EventInterface) (b Bootstrap, err error) {
	var global Global
	// 根据配置文件格式选择解析方法
	if isJson {
		global, err = JSON(file, event)
	} else {
		global, err = XML(file, event)
	}
	if err != nil {
		return
	}
	// 初始化RDS和数据库连接
	global.rdx, err = global.Rds()
	if err != nil {
		return
	}
	global.db, err = global.Db()
	if err != nil {
		return
	}
	// 注册全局英雄信息
	hero.Register(global)
	// 创建一个新的Iris应用实例
	app := iris.New()
	// 配置日志
	logbook, err := global.Log(global.Service.Log.File, global.Service.Log.Console)
	if err != nil {
		return
	}
	app.Logger().SetOutput(logbook).SetLevel(global.Service.Log.Level)
	return Bootstrap{
		Global: global,
		app:    app,
	}, nil
}

// Handle 方法用于处理指定的路由。
// 参数 route 是一个函数，它接受一个 iris.Application 的指针作为参数，
// 该函数通常用来配置和初始化 Iris 应用的路由。
// 返回值是 Bootstrap 类型的实例，实现了方法链的设计模式。
func (b Bootstrap) Handle(route func(app *iris.Application)) Bootstrap {
	// 调用传入的路由函数，传入当前 Bootstrap 实例的 app 字段。
	route(b.app)
	// 返回当前 Bootstrap 实例，支持方法链式调用。
	return b
}

// Set 是一个方法，它接收一个处理函数并使用该函数来处理当前的iris应用程序和全局配置。
// 这个方法允许开发者在应用程序启动过程中自定义一些逻辑。
// 参数：
//   handle - 一个函数，它接收一个iris应用程序实例和全局配置作为参数，用于执行自定义逻辑。
// 返回值：
//   Bootstrap - 返回当前的Bootstrap实例，允许链式调用。
func (b Bootstrap) Set(handle func(app *iris.Application, global Global)) Bootstrap {
	// 使用传入的handle函数处理当前的iris应用程序和全局配置
	handle(b.app, b.Global)
	return b
}

// Mvc 方法用于创建一个新的 mvc 应用，并注册全局中间件，然后允许调用者进一步处理和配置 mvc 应用。
// 参数：
//   handle - 一个函数，接收一个 mvc.Application 的指针，用于执行额外的配置或初始化。
// 返回值：
//   Bootstrap - 返回 Bootstrap 实例本身，支持链式调用。
func (b Bootstrap) Mvc(handle func(app *mvc.Application)) Bootstrap {
	b.m = mvc.New(b.app)   // 创建一个新的 mvc 应用实例
	b.m.Register(b.Global) // 注册全局中间件
	handle(b.m)            // 允许调用者进一步配置 mvc 应用
	return b               // 支持链式调用
}

// View 方法用于设置视图引擎并注册函数到模板引擎中。
// methods 参数是一个映射，其中键是函数在模板中的调用名称，值是对应的函数本身。
// 返回值是 Bootstrap 结构体，允许链式调用。
func (b Bootstrap) View(methods map[string]interface{}) Bootstrap {
	// 初始化 HTML 视图，并设置模板目录和扩展名
	view := iris.HTML(b.Global.Service.Template.Dir, b.Global.Service.Template.Ext)
	// 设置模板的左右定界符
	view.Delims(b.Global.Service.Template.Delimit[0], b.Global.Service.Template.Delimit[1])
	// 根据环境变量决定是否启用模板热加载
	view.Reload(strings.EqualFold(os.Getenv("ENV"), "development"))
	// 遍历 methods 映射，将所有方法注册到模板引擎
	for name, method := range methods {
		view.AddFunc(name, method)
	}
	// 注册视图引擎到 Iris 应用
	b.app.RegisterView(view)
	return b
}

// Run 启动应用程序。
// 该函数遍历全局服务配置中的资源目录，并根据存在与否处理静态资源目录和favicon。
// 然后，它在指定的主机和端口上启动应用程序。
func (b Bootstrap) Run(config iris.Configuration) {
	// 遍历配置的资源，如果目录存在，则将目录绑定到相应的URL上
	for _, resource := range b.Global.Service.Resources {
		if fsutil.PathExists(resource.Dir) {
			b.app.HandleDir(resource.Url, resource.Dir)
		}
	}
	// 检查Favicon文件是否存在，并设置
	if fsutil.FileExist(b.Global.Service.Favicon) {
		b.app.Favicon(b.Global.Service.Favicon)
	}
	config.DisableStartupLog = strings.EqualFold(os.Getenv("ENV"), "development")
	config.PostMaxMemory = b.Global.Service.Upload.Maximum << 20
	config.TimeFormat = "2006-01-02 15:04:05"
	config.Charset = "UTF-8"
	// 启动应用监听指定的主机和端口
	err := b.app.Run(
		iris.Addr(fmt.Sprintf("%s:%d", b.Global.Service.Host, b.Global.Service.Port)),
		iris.WithConfiguration(config),
	)
	if err != nil {
		panic(err) // 如果启动过程中出现错误，则抛出panic
	}
}
