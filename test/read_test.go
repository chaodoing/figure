package test

import (
	`encoding/json`
	`testing`
	
	`github.com/chaodoing/figure/app`
	`github.com/chaodoing/figure/toolkit`
)

func TestRead(t *testing.T) {
	var data app.Global
	err := toolkit.ReadXML("./index.xml", &data)
	if err != nil {
		t.Error(err)
	}
	s, e := json.MarshalIndent(data, "", "\t")
	if e != nil {
		t.Error(e)
	}
	t.Log(string(s))
}
