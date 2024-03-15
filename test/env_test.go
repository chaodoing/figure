package test

import (
	`os`
	`testing`
	
	`github.com/chaodoing/figure/app`
)

func TestEnv(t *testing.T) {
	_ = os.Setenv("DIR", "/Users/superman/Server/src/github.com/chaodoing/figure")
	err := app.GlobalDefault().MakeEnv()
	if err != nil {
		t.Error(err)
	}
	t.Log("Success")
}
