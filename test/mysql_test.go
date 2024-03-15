package test

import (
	`testing`
	
	`github.com/chaodoing/figure/app`
)

func TestMySQL(t *testing.T) {
	global := app.GlobalDefault()
	global.MySQL.Name = "pluvia"
	global.MySQL.Password = "123.com"
	global.Redis.Auth = "123.com"
	rdx, err := global.Rds()
	if err != nil {
		t.Error(err)
	}
	t.Log(rdx.Get("email").Result())
	// t.Log(rdx.Exists("email").Result())
	// t.Log(rdx.Set("email", "chaodoing@hotmail.com", 0).Result())
	// db, err := global.Db()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// var data map[string]interface{}
	// err = db.Table("administrator").Find(&data).Error
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Log(data)
}
