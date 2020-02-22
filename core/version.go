package core

import "fmt"

const (
	// MelodyVersion export the project version.
	MelodyVersion = "0.0.1"
	// MelodyHeaderKey
	MelodyHeaderKey = "X-Melody"
	//
)

var (
	// MelodyUserAgent setted to backend
	MelodyUserAgent   = fmt.Sprintf("Melody Version %s", MelodyVersion)
	MelodyHeaderValue = fmt.Sprintf("Version %s", MelodyVersion)

)
