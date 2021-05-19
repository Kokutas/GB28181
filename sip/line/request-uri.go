package line

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RequestUri struct {
	schema    string                 //schema
	user      string                 // user
	host      string                 // host
	port      uint16                 // port
	extension map[string]interface{} // extension
}

func (requestUri *RequestUri) SetSchema(schema string) {
	requestUri.schema = schema
}
func (requestUri *RequestUri) GetSchema() string {
	return requestUri.schema
}
func (requestUri *RequestUri) SetUser(user string) {
	requestUri.user = user
}
func (requestUri *RequestUri) GetUser() string {
	return requestUri.user
}
func (requestUri *RequestUri) SetHost(host string) {
	requestUri.host = host
}
func (requestUri *RequestUri) GetHost() string {
	return requestUri.host
}
func (requestUri *RequestUri) SetPort(port uint16) {
	requestUri.port = port
}
func (requestUri *RequestUri) GetPort() uint16 {
	return requestUri.port
}
func (requestUri *RequestUri) SetExtension(extension map[string]interface{}) {
	requestUri.extension = extension
}
func (requestUri *RequestUri) GetExtension() map[string]interface{} {
	return requestUri.extension
}
func NewRequestUri(schema, user, host string, port uint16, extension map[string]interface{}) *RequestUri {
	return &RequestUri{
		schema:    schema,
		user:      user,
		host:      host,
		port:      port,
		extension: extension,
	}
}

func (requestUri *RequestUri) Raw() (string, error) {
	result := ""
	if err := requestUri.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("%s:%s@%s", strings.ToLower(requestUri.schema), requestUri.user, requestUri.host)
	if requestUri.port > 0 {
		result += fmt.Sprintf(":%d", requestUri.port)
	}
	if requestUri.extension != nil {
		extensions := ""
		for k, v := range requestUri.extension {
			if len(strings.TrimSpace(fmt.Sprintf("%v", v))) == 0 {
				extensions += fmt.Sprintf(";%s", k)
			} else {
				extensions += fmt.Sprintf(";%s=%v", k, v)
			}
		}
		if len(strings.TrimSpace(extensions)) > 0 {
			result += extensions
		}
	}
	return result, nil
}
func (requestUri *RequestUri) Parse(raw string) error {
	if reflect.DeepEqual(nil, requestUri) {
		return errors.New("request-uri caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	schemaAndUserAndHostAndPortRegexp := regexp.MustCompile(`(?i)(sip):\d{20}@.*`)
	if schemaAndUserAndHostAndPortRegexp.MatchString(raw) {
		schemaAndUserAndHostAndPortStr := schemaAndUserAndHostAndPortRegexp.FindString(raw)
		schemaRegexp := regexp.MustCompile(`(?i)(sip):`)
		if schemaRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			requestUri.schema = strings.ToLower(regexp.MustCompile(`:`).ReplaceAllString(schemaRegexp.FindString(schemaAndUserAndHostAndPortStr), ""))
			schemaAndUserAndHostAndPortStr = schemaRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		userRegexp := regexp.MustCompile(`\d{20}@`)
		if userRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			requestUri.user = regexp.MustCompile(`@`).ReplaceAllString(userRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			schemaAndUserAndHostAndPortStr = userRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		portRegexp := regexp.MustCompile(`:\d+`)
		if portRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			portStr := regexp.MustCompile(`:`).ReplaceAllString(portRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return err
			}
			requestUri.port = uint16(port)
			schemaAndUserAndHostAndPortStr = portRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		schemaAndUserAndHostAndPortStr = strings.TrimPrefix(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimSuffix(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimPrefix(schemaAndUserAndHostAndPortStr, ";")
		schemaAndUserAndHostAndPortStr = strings.TrimSuffix(schemaAndUserAndHostAndPortStr, ";")
		if len(strings.TrimSpace(schemaAndUserAndHostAndPortStr)) > 0 && strings.Contains(schemaAndUserAndHostAndPortStr, ";") {
			parametersRegexp := regexp.MustCompile(`;.*`)
			m := make(map[string]interface{})
			extensions := strings.Split(parametersRegexp.FindString(schemaAndUserAndHostAndPortStr), ";")
			for _, v := range extensions {
				if len(strings.TrimSpace(v)) == 0 {
					continue
				}
				if strings.Contains(v, "=") {
					vs := strings.Split(v, "=")
					if len(vs) > 1 {
						m[vs[0]] = vs[1]
					} else {
						m[vs[0]] = ""
					}
				} else {
					m[v] = ""
				}
			}
			if len(m) > 0 {
				requestUri.extension = m
			}
			schemaAndUserAndHostAndPortStr = parametersRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		if len(strings.TrimSpace(schemaAndUserAndHostAndPortStr)) > 0 {
			requestUri.host = schemaAndUserAndHostAndPortStr
		}
	}
	return requestUri.Validator()
}

func (requestUri *RequestUri) Validator() error {
	if reflect.DeepEqual(nil, requestUri) {
		return errors.New("request-uri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(requestUri.schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(requestUri.schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if len(strings.TrimSpace(requestUri.user)) == 0 {
		return errors.New("the user field is not allowed to be empty")
	}
	if !regexp.MustCompile(`\d{20}`).MatchString(requestUri.user) {
		return errors.New("the user field must be 20 digits")
	}
	if len(strings.TrimSpace(requestUri.host)) == 0 {
		return errors.New("the host field is not allowed to be empty")
	}
	if ip := net.ParseIP(requestUri.host); ip != nil {
		if requestUri.port == 0 {
			return errors.New("the host field gives the ip address, the port must be given")
		}
	} else {
		if !reflect.DeepEqual(requestUri.host, requestUri.user[:10]) {
			return errors.New("the realm given by the host field must be consistent with the first 10 digits of the user")
		}
	}
	return nil
}

func (requestUri *RequestUri) String() string {
	result := ""
	result += fmt.Sprintf("%s:%s@%s", strings.ToLower(requestUri.schema), requestUri.user, requestUri.host)
	if requestUri.port > 0 {
		result += fmt.Sprintf(":%d", requestUri.port)
	}
	if requestUri.extension != nil {
		extensions := ""
		for k, v := range requestUri.extension {
			if len(strings.TrimSpace(fmt.Sprintf("%v", v))) == 0 {
				extensions += fmt.Sprintf(";%s", k)
			} else {
				extensions += fmt.Sprintf(";%s=%v", k, v)
			}
		}

		if len(strings.TrimSpace(extensions)) > 0 {
			result += extensions
		}
	}
	return result
}
