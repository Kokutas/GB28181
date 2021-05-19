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
	Method      string `json:"Method"`
	*RequestUri `json:"Request-URI"`
	Schema      string  `json:"Schema"`
	Version     float64 `json:"Version"`
}

func (requestLine *RequestLine) SetMethod(method string) {
	requestLine.Method = method
}
func (requestLine *RequestLine) GetMethod() string {
	return requestLine.Method
}
func (requestLine *RequestLine) SetRequestUri(uri *RequestUri) {
	requestLine.RequestUri = uri
}
func (requestLine *RequestLine) GetRequestUri() *RequestUri {
	if requestLine.RequestUri != nil {
		return requestLine.RequestUri
	}
	return nil
}
func (requestLine *RequestLine) SetSchema(schema string) {
	requestLine.Schema = schema
}
func (requestLine *RequestLine) GetSchema() string {
	return requestLine.Schema
}
func (requestLine *RequestLine) SetVersion(version float64) {
	requestLine.Version = version
}
func (requestLine *RequestLine) GetVersion() float64 {
	return requestLine.Version
}

func NewRequestLine(method string, uri *RequestUri, schema string, version float64) *RequestLine {
	return &RequestLine{
		Method:     method,
		RequestUri: uri,
		Schema:     schema,
		Version:    version,
	}
}
func (requestLine *RequestLine) Raw() (string, error) {
	result := ""
	if err := requestLine.Validator(); err != nil {
		return result, err
	}
	requestUri, err := requestLine.RequestUri.Raw()
	if err != nil {
		return result, fmt.Errorf("request-uri error : %s", err.Error())
	}
	result += fmt.Sprintf("%s %s %s/%1.1f", strings.ToUpper(requestLine.Method), requestUri, strings.ToUpper(requestLine.Schema), requestLine.Version)
	result += "\r\n"
	return result, nil
}

func (requestLine *RequestLine) Parse(raw string) error {
	if reflect.DeepEqual(nil, requestLine) {
		return errors.New("request-line caller is not allowed to be nil")
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
	requestLine.Method = strings.ToUpper(strings.TrimSpace(methodsRegexp.FindString(raw)))
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
	requestLine.Schema = strings.ToUpper(schemaRegexp.FindString(schemaAndVersionRegexp.FindString(raw)))
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
	requestLine.Version = version
	raw = schemaAndVersionRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the request-uri data cannot be parsed")
	}
	requestLine.RequestUri = new(RequestUri)

	if err := requestLine.RequestUri.Parse(raw); err != nil {
		return fmt.Errorf("request-uri parse error : %s", err.Error())
	}
	return requestLine.Validator()
}

func (requestLine *RequestLine) Validator() error {
	if reflect.DeepEqual(nil, requestLine) {
		return errors.New("request-line caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(requestLine.Method)) == 0 {
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
	if !methodsRegexp.MatchString(requestLine.Method) {
		return errors.New("the value of the method field cannot be matched")
	}
	if err := requestLine.RequestUri.Validator(); err != nil {
		return fmt.Errorf("request-line validator error : %s", err.Error())
	}
	if len(strings.TrimSpace(requestLine.Schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(requestLine.Schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if requestLine.Version != 2.0 {
		return errors.New("the value of the version field must be 2.0")
	}
	return nil
}
