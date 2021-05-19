package line

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kokutas/gb28181/sip/lib"
)

type StatusLine struct {
	Schema       string  `json:"Schema"`
	Version      float64 `json:"Version"`
	StatusCode   int     `json:"Status-Code"`
	ReasonPhrase string  `json:"Reason-Phrase"`
}

func (statusLine *StatusLine) SetSchema(schema string) {
	statusLine.Schema = schema
}
func (statusLine *StatusLine) GetSchema() string {
	return statusLine.Schema
}
func (statusLine *StatusLine) SetVersion(version float64) {
	statusLine.Version = version
}
func (statusLine *StatusLine) GetVersion() float64 {
	return statusLine.Version
}
func (statusLine *StatusLine) SetStatusCode(statusCode int) {
	statusLine.StatusCode = statusCode
}
func (statusLine *StatusLine) GetStatusCode() int {
	return statusLine.StatusCode
}
func (statusLine *StatusLine) SetReasonPhrase(reasonPhrase string) {
	statusLine.ReasonPhrase = reasonPhrase
}
func (statusLine *StatusLine) GetReasonPhrase() string {
	return statusLine.ReasonPhrase
}

func NewStatusLine(schema string, version float64, statusCode int, reasonPhrase string) *StatusLine {
	return &StatusLine{
		Schema:       schema,
		Version:      version,
		StatusCode:   statusCode,
		ReasonPhrase: reasonPhrase,
	}
}

func (statusLine *StatusLine) Raw() (string, error) {
	result := ""
	if err := statusLine.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("%s/%1.1f %d %s", strings.ToUpper(statusLine.Schema), statusLine.Version, statusLine.StatusCode, statusLine.ReasonPhrase)
	result += "\r\n"
	return result, nil
}
func (statusLine *StatusLine) Parse(raw string) error {
	if reflect.DeepEqual(nil, statusLine) {
		return errors.New("status-line caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// schema/version regexp
	schemaAndVersionRegexp := regexp.MustCompile(`(?i)(sip)/2\.0 `)
	if !schemaAndVersionRegexp.MatchString(raw) {
		return errors.New("the values of the schema and version fields cannot match")
	}
	// schema regexp
	schemaRegexp := regexp.MustCompile(`(?i)(sip)`)
	if !schemaRegexp.MatchString(schemaAndVersionRegexp.FindString(raw)) {
		return errors.New("the values of the schema field cannot match")
	}
	statusLine.Schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionRegexp.FindString(raw)))
	// version regexp
	versionRegexp := regexp.MustCompile(`2\.0`)
	if !versionRegexp.MatchString(raw) {
		return errors.New("the values of the version field cannot match")
	}
	versionStr := versionRegexp.FindString(schemaAndVersionRegexp.FindString(raw))
	version, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		return err
	}
	statusLine.Version = version
	raw = schemaAndVersionRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// status-code regexp
	statusCodeRegexp := regexp.MustCompile(`\d+`)
	if !statusCodeRegexp.MatchString(raw) {
		return errors.New("the values of the status-code field cannot parse")
	}
	statusCodeStr := statusCodeRegexp.FindString(raw)
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return err
	}
	statusLine.StatusCode = statusCode
	raw = statusCodeRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the reason-phrase data cannot be parsed")
	}
	statusLine.ReasonPhrase = raw
	return statusLine.Validator()
}
func (statusLine *StatusLine) Validator() error {
	if reflect.DeepEqual(nil, statusLine) {
		return errors.New("status-line caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(statusLine.Schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(statusLine.Schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if statusLine.Version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	if v1, ok1 := lib.Informational[statusLine.StatusCode]; ok1 {
		if !regexp.MustCompile(`(?i)` + v1).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v2, ok2 := lib.Success[statusLine.StatusCode]; ok2 {
		if !regexp.MustCompile(`(?i)` + v2).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v3, ok3 := lib.Redirection[statusLine.StatusCode]; ok3 {
		if !regexp.MustCompile(`(?i)` + v3).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v4, ok4 := lib.ClientError[statusLine.StatusCode]; ok4 {
		if !regexp.MustCompile(`(?i)` + v4).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v5, ok5 := lib.ServerError[statusLine.StatusCode]; ok5 {
		if !regexp.MustCompile(`(?i)` + v5).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v6, ok6 := lib.GlobalFailure[statusLine.StatusCode]; ok6 {
		if !regexp.MustCompile(`(?i)` + v6).MatchString(statusLine.ReasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else {
		return errors.New("the value of the status-code and reason-phrase fields is not match")
	}
	return nil
}
