package test

import (
	`testing`
	
	`github.com/chaodoing/figure/app`
	`github.com/chaodoing/figure/toolkit`
)

func TestEnv(t *testing.T) {
	var data = app.Global{
		Service: &app.Service{
			Favicon:     "${DIR}/resources/favicon.ico",
			Port:        8080,
			Host:        "127.0.0.1",
			CrossDomain: true,
			Log: &app.Logger{
				Console: true,
				File:    "${DIR}/logs/app-%Y-%m-%d.log",
				Level:   "debug",
			},
			Resources: []app.Resource{
				{Url: "/upload", Dir: "${DIR}/resources/upload"},
				{Url: "/static", Dir: "${DIR}/resources/static"},
			},
			Template: &app.Template{
				Delimit: []string{"{{", "}}"},
				Dir:     "${DIR}/resources/template",
				Ext:     ".html",
			},
		},
		MySQL: &app.MySQL{
			Host:     "127.0.0.1",
			Port:     3306,
			Name:     "pluvia",
			Username: "root",
			Password: "123.com",
			Charset:  "utf8mb4",
			Logger: &app.Logger{
				Console: true,
				File:    "${DIR}/logs/mysql-%Y-%m-%d.log",
				Level:   "info",
			},
		},
		Redis: &app.Redis{
			Host: "127.0.0.1",
			Port: 6379,
			Db:   0,
			Auth: "123.com",
			TTL:  648000,
		},
	}
	err := toolkit.SaveXML(data, "./index.xml")
	if err != nil {
		t.Error(err)
	}
	t.Log("Success")
}
