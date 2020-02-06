package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Parser returns a ServiceConfig struct according to the providing configFile name.
type Parser interface {
	Parse(configFile string) (ServiceConfig, error)
}

// NewParseError returns a new ParseError
func NewParseError(err error, configFile string, offset int) *ParseError {
	b, _ := ioutil.ReadFile(configFile)
	row, col := getErrorRowCol(b, offset)
	return &ParseError{
		ConfigFile: configFile,
		Err:        err,
		Offset:     offset,
		Row:        row,
		Col:        col,
	}
}

//CheckErr returns error when parse config file
func CheckErr(err error, configFile string) error {
	switch e := err.(type) {
	case *json.SyntaxError:
		return NewParseError(err, configFile, int(e.Offset))
	case *json.UnmarshalTypeError:
		return NewParseError(err, configFile, int(e.Offset))
	case *os.PathError:
		return fmt.Errorf(
			"'%s' (%s): %s",
			configFile,
			e.Op,
			e.Err.Error(),
		)
	default:
		return fmt.Errorf("'%s': %v", configFile, err)
	}
}

func getErrorRowCol(source []byte, offset int) (row, col int) {
	for i := 0; i < offset; i++ {
		v := source[i]
		if v == '\r' {
			continue
		}
		if v == '\n' {
			col = 0
			row++
			continue
		}
		col++
	}
	return
}

// ParseError is an error containing details regarding the row and column where
// an parse error occurred
type ParseError struct {
	ConfigFile string
	Offset     int
	Row        int
	Col        int
	Err        error
}

// Error returns the error message for the ParseError
func (p *ParseError) Error() string {
	return fmt.Sprintf(
		"'%s': %v, offset: %v, row: %v, col: %v",
		p.ConfigFile,
		p.Err.Error(),
		p.Offset,
		p.Row,
		p.Col,
	)
}
