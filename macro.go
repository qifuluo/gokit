package gokit

//chan大小
const (
	CHANSIZE_LOG      = 1024 * 10      //日志
)

//日志级别
const (
	_ = iota
	LV_SYS
	LV_ERROR
	LV_WARN
	LV_INFO
	LV_DEBUG
)
