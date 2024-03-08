package app

import (
	`encoding/json`
	`encoding/xml`
	`os`
)

// XML 加载配置文件
func XML(env string) (data *Global, err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(env))
	if err != nil {
		return
	}
	err = xml.Unmarshal(content, &data)
	if err != nil {
		return
	}
	c := data.LoadEnv()
	if c != nil {
		data = c
	}
	return
}

// JSON 加载配置文件
func JSON(env string) (data *Global, err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(env))
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}
	c := data.LoadEnv()
	if c != nil {
		data = c
	}
	return
}
