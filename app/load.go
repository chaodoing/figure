package app

import (
	`encoding/json`
	`encoding/xml`
	`os`
)

// XML 加载配置文件
func XML(env string, event EventInterface) (data Global, err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(env))
	if err != nil {
		return
	}
	err = xml.Unmarshal(content, &data)
	if err != nil {
		return
	}
	data = data.LoadEnv(event)
	return
}

// JSON 加载配置文件
func JSON(env string, event EventInterface) (data Global, err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(env))
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return
	}
	data = data.LoadEnv(event)
	return
}
