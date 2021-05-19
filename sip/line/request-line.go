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

type RequestLine struct {
	method  string      // method
	reqUri  *RequestUri // request-uri
	schema  string      // schema
	version float64     // version
}

func (requestLine *RequestLine) SetMethod(method string) {
	requestLine.method = method
}
func (requestLine *RequestLine) GetMethod() string {
	return requestLine.method
}
func (requestLine *RequestLine) SetReqUri(uri *RequestUri) {
	requestLine.reqUri = uri
}
func (requestLine *RequestLine) GetReqUri() *RequestUri {
	return requestLine.reqUri
}
func (requestLine *RequestLine) SetSchema(schema string) {
	requestLine.schema = schema
}
func (requestLine *RequestLine) GetSchema() string {
	return requestLine.schema
}
func (requestLine *RequestLine) SetVersion(version float64) {
	requestLine.version = version
}
func (requestLine *RequestLine) GetVersion() float64 {
	return requestLine.version
}

func NewRequestLine(method string, reqUri *RequestUri, schema string, version float64) *RequestLine {
	return &RequestLine{
		method:  method,
		reqUri:  reqUri,
		schema:  schema,
		version: version,
	}
}
func (requestLine *RequestLine) Raw() (string, error) {
	result := ""
	if err := requestLine.Validator(); err != nil {
		return result, err
	}
	requestUri, err := requestLine.reqUri.Raw()
	if err != nil {
		return result, fmt.Errorf("request-uri error : %s", err.Error())
	}
	result += fmt.Sprintf("%s %s %s/%1.1f", strings.ToUpper(requestLine.method), requestUri, strings.ToUpper(requestLine.schema), requestLine.version)
	result += "\r\n"
	return result, nil
}

func (requestLine *RequestLine) Parse(raw string) error {
	if reflect.DeepEqual(nil, requestLine) {
		return errors.New("request-line caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range lib.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ") "
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if !methodsRegexp.MatchString(raw) {
		return errors.New("the value of the method field cannot be matched")
	}
	requestLine.method = strings.ToUpper(strings.TrimSpace(methodsRegexp.FindString(raw)))
	raw = methodsRegexp.ReplaceAllString(raw, "")

	// schema/version regexp
	schemaAndVersionRegexp := regexp.MustCompile(`(?i)(sip)/2\.0$`)
	if !schemaAndVersionRegexp.MatchString(raw) {
		return errors.New("the values of the schema and version fields cannot match")
	}
	// schema regexp
	schemaRegexp := regexp.MustCompile(`(?i)(sip)`)
	if !schemaRegexp.MatchString(schemaAndVersionRegexp.FindString(raw)) {
		return errors.New("the value of the schema field cannot match")
	}
	requestLine.schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionRegexp.FindString(raw)))
	// version regexp
	versionRegexp := regexp.MustCompile(`2\.0`)
	if !versionRegexp.MatchString(schemaAndVersionRegexp.FindString(raw)) {
		return errors.New("the value of the version field cannot match")
	}
	versionStr := versionRegexp.FindString(schemaAndVersionRegexp.FindString(raw))
	version, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		return err
	}
	requestLine.version = version
	raw = schemaAndVersionRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the request-uri data cannot be parsed")
	}
	requestLine.reqUri = new(RequestUri)

	if err := requestLine.reqUri.Parse(raw); err != nil {
		return fmt.Errorf("request-uri parse error : %s", err.Error())
	}
	return requestLine.Validator()
}

func (requestLine *RequestLine) Validator() error {
	if reflect.DeepEqual(nil, requestLine) {
		return errors.New("request-line caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(requestLine.method)) == 0 {
		return errors.New("request-line method field is not allowed to be empty")
	}
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range lib.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")$"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if !methodsRegexp.MatchString(requestLine.method) {
		return errors.New("the value of the method field cannot be matched")
	}
	if err := requestLine.reqUri.Validator(); err != nil {
		return fmt.Errorf("request-line validator error : %s", err.Error())
	}
	if len(strings.TrimSpace(requestLine.schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(requestLine.schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if requestLine.version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	return nil
}

func (requestLine *RequestLine) String() string {
	result := ""
	if len(strings.TrimSpace(requestLine.method)) > 0 {
		result += strings.ToUpper(requestLine.method)
	}
	if !reflect.DeepEqual(nil, requestLine.reqUri) {
		if len(result) > 0 {
			result += fmt.Sprintf(" %s", requestLine.reqUri.String())
		} else {
			result += requestLine.reqUri.String()
		}
	}
	if len(strings.TrimSpace(requestLine.schema)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(" %s/%1.1f", strings.ToUpper(requestLine.schema), requestLine.version)
		} else {
			result += fmt.Sprintf("%s/%1.1f", strings.ToUpper(requestLine.schema), requestLine.version)
		}
	}
	return result
}
