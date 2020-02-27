package core

import "fmt"

const (
	// MelodyVersion export the project version.
	MelodyVersion = "1.0.0"
	// MelodyHeaderKey
	MelodyHeaderKey = "X-Melody"
	//
)

var (
	// MelodyUserAgent setted to backend
	MelodyUserAgent   = fmt.Sprintf("Melody Version %s", MelodyVersion)
	MelodyHeaderValue = fmt.Sprintf("Version %s", MelodyVersion)

)
