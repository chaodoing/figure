package test

import (
	`testing`
	
	`github.com/chaodoing/figure/app`
	`github.com/chaodoing/figure/toolkit`
	`github.com/gookit/goutil/fsutil`
	encoder `github.com/zwgblue/yaml-encoder`
)

func TestYaml(t *testing.T) {
	var data app.Global
	err := toolkit.ReadXML("./index.xml", &data)
	if err != nil {
		t.Error(err)
	}
	v, e := encoder.NewEncoder(data, encoder.WithComments(encoder.CommentsOnHead)).Encode()
	if e != nil {
		t.Error(e)
	}
	_, err = fsutil.PutContents("./env.yaml", v)
	if err != nil {
		t.Error(err)
	}
	t.Log("Success")
}
