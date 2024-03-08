package toolkit

import (
	`encoding/json`
	`encoding/xml`
	`os`
)

// ReadJSON 读取JSON
func ReadJSON(file string, data interface{}) (err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(file))
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	return nil
}

// SaveJSON 写入JSON
func SaveJSON(data interface{}, file string) error {
	xmlByte, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(os.ExpandEnv(file), xmlByte, 0666); err != nil {
		return err
	}
	return nil
}

// ReadXML 读取XML文件
func ReadXML(file string, data interface{}) (err error) {
	var content []byte
	content, err = os.ReadFile(os.ExpandEnv(file))
	if err != nil {
		return err
	}
	err = xml.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	return nil
}

// SaveXML 存储XML文件
func SaveXML(data interface{}, file string) error {
	xmlByte, err := xml.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	headerBytes := []byte(xml.Header)
	xmlData := append(headerBytes, xmlByte...)
	if err := os.WriteFile(os.ExpandEnv(file), xmlData, 0666); err != nil {
		return err
	}
	return nil
}
