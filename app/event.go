package app

type Event struct{}

// AppInit 应用初始化
func (e *Event) AppInit() {
	
}

// DbInit 数据库初始化
func (e *Event) DbInit() {
	
}

// RDSInit redis初始化
func (e *Event) RDSInit() {}

// HttpRun http服务启动
func (e *Event) HttpRun() {
	
}

// HttpEnd http服务结束
func (e *Event) HttpEnd() {
	
}
