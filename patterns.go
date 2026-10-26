package validatex

import (
	"regexp"
	"time"
)

var (
	//uuidPattern 定义UUID正则
	uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	//ulidPattern 定义ULID正则
	ulidPattern = regexp.MustCompile(`(?i)^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)
	//hexColorPattern 定义十六进制颜色正则
	hexColorPattern = regexp.MustCompile(`(?i)^#(?:[0-9a-f]{3}|[0-9a-f]{6}|[0-9a-f]{8})$`)
)

// dateLayouts 定义常用日期布局
var dateLayouts = [...]string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	time.RFC1123,
}
