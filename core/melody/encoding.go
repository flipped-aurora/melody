package melody

import (
	rss "melody/middleware/melody-rss"
	xml "melody/middleware/melody-xml"
)

// RegisterEncoders 注册额外的编码器
func RegisterEncoders() {
	xml.Register()
	rss.Register()
}
