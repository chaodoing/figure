package main

import (
	`os`
	
	`github.com/chaodoing/figure/app`
	`github.com/chaodoing/figure/models`
	`github.com/chaodoing/figure/o`
	`github.com/gookit/goutil/envutil`
	`github.com/kataras/iris/v12`
	`github.com/kataras/iris/v12/hero`
	`github.com/kataras/iris/v12/mvc`
)

func index(ctx iris.Context, global app.Global) {
	db, err := global.Db()
	if err != nil {
		o.O(ctx, o.Data{Code: 3306, Message: err.Error()})
		return
	}
	var data []models.Admin
	err = db.Find(&data).Error
	if err != nil {
		o.O(ctx, o.Data{Code: 1, Message: err.Error()})
		return
	}
	o.O(ctx, o.Data{Code: 0, Message: `success`, Data: data})
	return
}

var (
	ENV     = "development"
	APP     = "figure"
	VERSION = "v1.0.0"
)

func main() {
	envutil.SetEnvMap(map[string]string{
		"DIR":     envutil.Getenv("DIR", os.ExpandEnv("${PWD}")),
		"ENV":     envutil.Getenv("ENV", ENV),
		"APP":     envutil.Getenv("APP", APP),
		"VERSION": envutil.Getenv("VERSION", VERSION),
	})
	boot, err := app.New("./config/app.xml", false, nil)
	if err != nil {
		panic(err)
	}
	
	boot.Handle(func(app *iris.Application) {
		app.Get(`/index`, hero.Handler(index))
	}).Mvc(func(app *mvc.Application) {
		app.Party("/api").Handle(new(Controller)).Name = "api"
	}).Run(iris.DefaultConfiguration())
}
