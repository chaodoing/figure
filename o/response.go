package o

import (
	`bytes`
	_ `embed`
	`encoding/json`
	`encoding/xml`
	`html/template`
	
	`github.com/gookit/goutil/fsutil`
	`github.com/kataras/iris/v12`
)

//go:embed vscode/index.html
var vscode string

type (
	// Response 结构体封装了 iris.Context 上下文信息。
	Response struct {
		ctx iris.Context // iris.Context 提供了处理 HTTP 请求的方法和数据。
	}
	
	// Data 结构体用于封装 API 响应的基本信息，包括状态码、消息和数据。
	Data struct {
		XMLName xml.Name    `json:"-" xml:"root" yaml:"-"`                                   // XML 名称，但在 JSON 和 YAML 中不使用。
		Code    int         `json:"code" xml:"code" yaml:"Code" comment:"响应状态码"`        // Code 响应状态码，用于表示请求处理的结果状态。
		Message string      `json:"message" xml:"message" yaml:"Message" comment:"响应消息"` // Message 响应消息，用于提供请求处理的详细信息。
		Data    interface{} `json:"data" xml:"data" yaml:"Data" comment:"响应数据"`          // Data 响应数据，实际返回给客户端的数据内容。
	}
	
	// Pagination 结构体用于封装分页信息，方便在分页场景下使用。
	Pagination struct {
		XMLName xml.Name    `json:"-" xml:"root" yaml:"-"`                                   // XML 名称，但在 JSON 和 YAML 中不使用。
		Page    int         `json:"page" xml:"page" yaml:"Page" comment:"当前页码"`          // Page 当前页码，表示当前请求的是哪一页的数据。
		Total   int         `json:"total" xml:"total" yaml:"Total" comment:"总条数"`         // Total 总条数，表示数据总共有多少条。
		Size    int         `json:"size" xml:"size" yaml:"Size" comment:"每页条数"`          // Size 每页条数，表示每页显示的数据数量。
		Code    int         `json:"code" xml:"code" yaml:"Code" comment:"响应状态码"`        // Code 响应状态码，用于表示分页请求处理的结果状态。
		Message string      `json:"message" xml:"message" yaml:"Message" comment:"响应消息"` // Message 响应消息，用于提供分页请求处理的详细信息。
		Data    interface{} `json:"data" xml:"data" yaml:"Data" comment:"响应数据"`          // Data 响应数据，实际返回给客户端的分页数据内容。
	}
)

func html(ctx iris.Context, data interface{}) (value string, err error) {
	tpl, err := template.New("json").Parse(vscode)
	if err != nil {
		return
	}
	js, err := json.Marshal(data)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, map[string]string{
		"Title": "JSON",
		"Json":  string(js),
		"Theme": ctx.URLParamDefault("theme", "vs"),
	})
	value = buf.String()
	return
}

func MD(ctx iris.Context, file string, data interface{}) {
	var value = fsutil.GetContents(file)
	tpl, err := template.New("markdown").Parse(string(value))
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	err = ctx.Markdown(buf.Bytes())
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	return
}

func O(ctx iris.Context, data interface{}) {
	var (
		err   error
		value string
	)
	value, err = html(ctx, data)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	ctx.Negotiation().HTML(value).JSON(data).XML(data).YAML(data).MsgPack(data).EncodingGzip().Charset("UTF-8")
	_, err = ctx.Negotiate(nil)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	return
}
