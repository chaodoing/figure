package main

import (
	`net/http`
	
	`github.com/chaodoing/figure/app`
	`github.com/chaodoing/figure/models`
	`github.com/chaodoing/figure/o`
	`github.com/kataras/iris/v12`
	`github.com/kataras/iris/v12/mvc`
)

type Controller struct{}

func (c *Controller) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/{id:uint64}", "Edit", func(ctx iris.Context) {
		ctx.Next()
	})
	b.Handle(http.MethodGet, "/{name:string}", "Name", func(ctx iris.Context) {
		ctx.Next()
	})
}

func (c *Controller) Name(ctx iris.Context) {
	err := ctx.Markdown([]byte(`# Name`))
	if err != nil {
		ctx.Application().Logger().Error(err)
		return
	}
}

func (c *Controller) Edit(ctx iris.Context, global app.Global) {
	var value o.Data
	db, err := global.Db()
	if err != nil {
		value = o.Data{Code: 3306, Message: err.Error()}
	}
	id := ctx.Params().GetUint64Default("id", 0)
	var data models.Admin
	err = db.Where("`id`=?", id).First(&data).Error
	if err != nil {
		value = o.Data{Code: 3306, Message: err.Error()}
	}
	value = o.Data{Code: 0, Message: "OK", Data: data}
	o.O(ctx, value)
}

func (c *Controller) Get(ctx iris.Context, global app.Global) {
	var value o.Data
	var data []models.Admin
	db, err := global.Db()
	if err != nil {
		value = o.Data{
			Code:    3306,
			Message: err.Error(),
		}
	}
	err = db.Table("admin").Find(&data).Error
	if err != nil {
		value = o.Data{
			Code:    1,
			Message: err.Error(),
		}
	}
	value = o.Data{
		Code:    0,
		Message: "OK",
		Data:    data,
	}
	o.O(ctx, value)
}
