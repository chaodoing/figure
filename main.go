package main

import (
	`github.com/kataras/iris/v12`
)

// type Counter struct {
// 	Service services.
// }
//
// func (c *Counter) BeforeActivation(b mvc.BeforeActivation) {
// 	b.Handle("GET", "/counter/{id:long}", "Editor")
// }
//
// func (c *Counter) HandleHTTPError(ctx iris.Context) {
// 	_ = ctx.JSON(iris.Map{
// 		"error":   ctx.GetStatusCode(),
// 		"message": "error",
// 	})
// }
//
// func (c *Counter) Editor(ctx iris.Context) mvc.View {
// 	return mvc.View{
// 		Name: "index.html",
// 	}
// }
//
// func (c *Counter) Get(ctx iris.Context) iris.Map {
// 	value := c.Service.Increment()
// 	return iris.Map{
// 		"value": value,
// 	}
// }

func main() {
	app := iris.Default()
	app.Any(`/`, func(ctx iris.Context) {
		_ = ctx.JSON(iris.Map{
			"error":   0,
			"message": "OK",
		})
	})
	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
