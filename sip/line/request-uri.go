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
	Schema    string                 `json:"Schema"`
	User      string                 `json:"User"`
	Host      string                 `json:"Host"`
	Port      uint16                 `json:"Port"`
	Extension map[string]interface{} `json:"Extension"`
}

func (uri *RequestUri) SetSchema(schema string) {
	uri.Schema = schema
}
func (uri *RequestUri) GetSchema() string {
	return uri.Schema
}
func (uri *RequestUri) SetUser(user string) {
	uri.User = user
}
func (uri *RequestUri) GetUser() string {
	return uri.User
}
func (uri *RequestUri) SetHost(host string) {
	uri.Host = host
}
func (uri *RequestUri) GetHost() string {
	return uri.Host
}
func (uri *RequestUri) SetPort(port uint16) {
	uri.Port = port
}
func (uri *RequestUri) GetPort() uint16 {
	return uri.Port
}
func (uri *RequestUri) SetExtension(extension map[string]interface{}) {
	uri.Extension = extension
}
func (uri *RequestUri) GetExtension() map[string]interface{} {
	if uri.Extension != nil {
		return uri.Extension
	}
	return nil
}
func NewRequestUri(schema, user, host string, port uint16, extension map[string]interface{}) *RequestUri {
	return &RequestUri{
		Schema:    schema,
		User:      user,
		Host:      host,
		Port:      port,
		Extension: extension,
	}
}

func (uri *RequestUri) Raw() (string, error) {
	result := ""
	if err := uri.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("%s:%s@%s", strings.ToLower(uri.Schema), uri.User, uri.Host)
	if uri.Port > 0 {
		result += fmt.Sprintf(":%d", uri.Port)
	}
	if uri.Extension != nil {
		extensions := ""
		for k, v := range uri.Extension {
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
func (uri *RequestUri) Parse(raw string) error {
	if reflect.DeepEqual(nil, uri) {
		return errors.New("request-uri caller is not allowed to be nil")
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
	schemaAndUserAndHostAndPortRegexp := regexp.MustCompile(`(?i)(sip):\d{20}@.*`)
	if schemaAndUserAndHostAndPortRegexp.MatchString(raw) {
		schemaAndUserAndHostAndPortStr := schemaAndUserAndHostAndPortRegexp.FindString(raw)
		schemaRegexp := regexp.MustCompile(`(?i)(sip):`)
		if schemaRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			uri.Schema = strings.ToLower(regexp.MustCompile(`:`).ReplaceAllString(schemaRegexp.FindString(schemaAndUserAndHostAndPortStr), ""))
			schemaAndUserAndHostAndPortStr = schemaRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		userRegexp := regexp.MustCompile(`\d{20}@`)
		if userRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			uri.User = regexp.MustCompile(`@`).ReplaceAllString(userRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			schemaAndUserAndHostAndPortStr = userRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		portRegexp := regexp.MustCompile(`:\d+`)
		if portRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			portStr := regexp.MustCompile(`:`).ReplaceAllString(portRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return err
			}
			uri.Port = uint16(port)
			schemaAndUserAndHostAndPortStr = portRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		schemaAndUserAndHostAndPortStr = strings.TrimLeft(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimRight(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimPrefix(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimSuffix(schemaAndUserAndHostAndPortStr, " ")
		schemaAndUserAndHostAndPortStr = strings.TrimLeft(schemaAndUserAndHostAndPortStr, ";")
		schemaAndUserAndHostAndPortStr = strings.TrimRight(schemaAndUserAndHostAndPortStr, ";")
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
				uri.Extension = m
			}
			schemaAndUserAndHostAndPortStr = parametersRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		if len(strings.TrimSpace(schemaAndUserAndHostAndPortStr)) > 0 {
			uri.Host = schemaAndUserAndHostAndPortStr
		}
	}
	return uri.Validator()

	// schemaRegexp := regexp.MustCompile(`(?i)(sip):`)
	// if schemaRegexp.MatchString(raw) {
	// 	uri.Schema = strings.ToLower(regexp.MustCompile(`:`).ReplaceAllString(schemaRegexp.FindString(raw), ""))
	// 	raw = schemaRegexp.ReplaceAllString(raw, "")
	// }
	// userRegexp := regexp.MustCompile(`\d{20}@`)
	// if userRegexp.MatchString(raw) {
	// 	uri.User = regexp.MustCompile(`@`).ReplaceAllString(userRegexp.FindString(raw), "")
	// 	raw = userRegexp.ReplaceAllString(raw, "")
	// }
	// portRegexp := regexp.MustCompile(`:\d+`)
	// if portRegexp.MatchString(raw) {
	// 	portStr := regexp.MustCompile(`:`).ReplaceAllString(portRegexp.FindString(raw), "")
	// 	port, err := strconv.Atoi(portStr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	uri.Port = uint16(port)
	// 	raw = portRegexp.ReplaceAllString(raw, "")
	// }
	// if len(strings.TrimSpace(raw)) > 0 {
	// 	uri.Host = raw
	// }
	// return uri.Validator()
}

func (uri *RequestUri) Validator() error {
	if reflect.DeepEqual(nil, uri) {
		return errors.New("request-uri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(uri.Schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(uri.Schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if len(strings.TrimSpace(uri.User)) == 0 {
		return errors.New("the user field is not allowed to be empty")
	}
	if !regexp.MustCompile(`\d{20}`).MatchString(uri.User) {
		return errors.New("the user field must be 20 digits")
	}
	if len(strings.TrimSpace(uri.Host)) == 0 {
		return errors.New("the host field is not allowed to be empty")
	}
	if ip := net.ParseIP(uri.Host); ip != nil {
		if uri.Port == 0 {
			return errors.New("the host field gives the ip address, the port must be given")
		}
	} else {
		if !reflect.DeepEqual(uri.Host, uri.User[:10]) {
			return errors.New("the realm given by the host field must be consistent with the first 10 digits of the user")
		}
	}
	return nil
}
