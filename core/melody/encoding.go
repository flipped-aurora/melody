package melody

import xml "melody/middleware/melody-xml"

func RegisterEncoders() {
	xml.Register()
	//TODO RSS register
}
