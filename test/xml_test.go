package test

import (
	`encoding/xml`
	`testing`
	
	`github.com/chaodoing/figure/app`
	`github.com/gookit/goutil/fsutil`
)

func TestXML(t *testing.T) {
	global := app.GlobalDefault()
	x, err := xml.Marshal(global)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = fsutil.PutContents("../config/app.xml", x)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success")
}
