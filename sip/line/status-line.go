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
	schema       string  // schema
	version      float64 // version
	statusCode   int     // Status-Code
	reasonPhrase string  // Reason-Phrase
}

func (statusLine *StatusLine) SetSchema(schema string) {
	statusLine.schema = schema
}
func (statusLine *StatusLine) GetSchema() string {
	return statusLine.schema
}
func (statusLine *StatusLine) SetVersion(version float64) {
	statusLine.version = version
}
func (statusLine *StatusLine) GetVersion() float64 {
	return statusLine.version
}
func (statusLine *StatusLine) SetStatusCode(statusCode int) {
	statusLine.statusCode = statusCode
}
func (statusLine *StatusLine) GetStatusCode() int {
	return statusLine.statusCode
}
func (statusLine *StatusLine) SetReasonPhrase(reasonPhrase string) {
	statusLine.reasonPhrase = reasonPhrase
}
func (statusLine *StatusLine) GetReasonPhrase() string {
	return statusLine.reasonPhrase
}

func NewStatusLine(schema string, version float64, statusCode int, reasonPhrase string) *StatusLine {
	return &StatusLine{
		schema:       schema,
		version:      version,
		statusCode:   statusCode,
		reasonPhrase: reasonPhrase,
	}
}

func (statusLine *StatusLine) Raw() (string, error) {
	result := ""
	if err := statusLine.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("%s/%1.1f %d %s", strings.ToUpper(statusLine.schema), statusLine.version, statusLine.statusCode, statusLine.reasonPhrase)
	result += "\r\n"
	return result, nil
}
func (statusLine *StatusLine) Parse(raw string) error {
	if reflect.DeepEqual(nil, statusLine) {
		return errors.New("status-line caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
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
	statusLine.schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionRegexp.FindString(raw)))
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
	statusLine.version = version
	raw = schemaAndVersionRegexp.ReplaceAllString(raw, "")
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
	statusLine.statusCode = statusCode
	raw = statusCodeRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the reason-phrase data cannot be parsed")
	}
	statusLine.reasonPhrase = raw
	return statusLine.Validator()
}
func (statusLine *StatusLine) Validator() error {
	if reflect.DeepEqual(nil, statusLine) {
		return errors.New("status-line caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(statusLine.schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(statusLine.schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if statusLine.version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	if v1, ok1 := lib.Informational[statusLine.statusCode]; ok1 {
		if !regexp.MustCompile(`(?i)` + v1).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v2, ok2 := lib.Success[statusLine.statusCode]; ok2 {
		if !regexp.MustCompile(`(?i)` + v2).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v3, ok3 := lib.Redirection[statusLine.statusCode]; ok3 {
		if !regexp.MustCompile(`(?i)` + v3).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v4, ok4 := lib.ClientError[statusLine.statusCode]; ok4 {
		if !regexp.MustCompile(`(?i)` + v4).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v5, ok5 := lib.ServerError[statusLine.statusCode]; ok5 {
		if !regexp.MustCompile(`(?i)` + v5).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else if v6, ok6 := lib.GlobalFailure[statusLine.statusCode]; ok6 {
		if !regexp.MustCompile(`(?i)` + v6).MatchString(statusLine.reasonPhrase) {
			return errors.New("the value of the reason-phrase field is not match")
		}
	} else {
		return errors.New("the value of the status-code and reason-phrase fields is not match")
	}
	return nil
}
func (statusLine *StatusLine) String() string {
	result := ""
	if len(strings.TrimSpace(statusLine.schema)) > 0 {
		result += fmt.Sprintf("%s/%1.1f", strings.ToUpper(statusLine.schema), statusLine.version)
	}
	if len(result) > 0 {
		result += fmt.Sprintf(" %d", statusLine.statusCode)
	} else {
		result += fmt.Sprintf("%d", statusLine.statusCode)
	}
	if len(strings.TrimSpace(statusLine.reasonPhrase)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(" %s", statusLine.reasonPhrase)
		} else {
			result += statusLine.reasonPhrase
		}
	}
	return result
}
