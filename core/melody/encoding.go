package melody

import (
	rss "melody/middleware/melody-rss"
	xml "melody/middleware/melody-xml"
)

func RegisterEncoders() {
	xml.Register()
	rss.Register()
}
