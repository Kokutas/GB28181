package header

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Uri struct {
	schema    string                 // schema
	user      string                 // user
	host      string                 // host
	port      uint16                 // port
	extension map[string]interface{} // extension
}

func (uri *Uri) SetSchema(schema string) {
	uri.schema = schema
}
func (uri *Uri) GetSchema() string {
	return uri.schema
}
func (uri *Uri) SetUser(user string) {
	uri.user = user
}
func (uri *Uri) GetUser() string {
	return uri.user
}
func (uri *Uri) SetHost(host string) {
	uri.host = host
}
func (uri *Uri) GetHost() string {
	return uri.host
}
func (uri *Uri) SetPort(port uint16) {
	uri.port = port
}
func (uri *Uri) GetPort() uint16 {
	return uri.port
}
func (uri *Uri) SetExtension(extension map[string]interface{}) {
	uri.extension = extension
}
func (uri *Uri) GetExtension() map[string]interface{} {
	return uri.extension
}

func NewUri(schema, user, host string, port uint16, extension map[string]interface{}) *Uri {
	return &Uri{
		schema:    schema,
		user:      user,
		host:      host,
		port:      port,
		extension: extension,
	}
}

func (uri *Uri) Raw() (string, error) {
	result := ""
	if err := uri.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("%s:%s@%s", strings.ToLower(uri.schema), uri.user, uri.host)
	if uri.port > 0 {
		result += fmt.Sprintf(":%d", uri.port)
	}
	if uri.extension != nil {
		extensions := ""
		for k, v := range uri.extension {
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
func (uri *Uri) Parse(raw string) error {
	if reflect.DeepEqual(nil, uri) {
		return errors.New("uri caller is not allowed to be nil")
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
			uri.schema = strings.ToLower(regexp.MustCompile(`:`).ReplaceAllString(schemaRegexp.FindString(schemaAndUserAndHostAndPortStr), ""))
			schemaAndUserAndHostAndPortStr = schemaRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		userRegexp := regexp.MustCompile(`\d{20}@`)
		if userRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			uri.user = regexp.MustCompile(`@`).ReplaceAllString(userRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			schemaAndUserAndHostAndPortStr = userRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		portRegexp := regexp.MustCompile(`:\d+`)
		if portRegexp.MatchString(schemaAndUserAndHostAndPortStr) {
			portStr := regexp.MustCompile(`:`).ReplaceAllString(portRegexp.FindString(schemaAndUserAndHostAndPortStr), "")
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return err
			}
			uri.port = uint16(port)
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
				uri.extension = m
			}
			schemaAndUserAndHostAndPortStr = parametersRegexp.ReplaceAllString(schemaAndUserAndHostAndPortStr, "")
		}
		if len(strings.TrimSpace(schemaAndUserAndHostAndPortStr)) > 0 {
			uri.host = schemaAndUserAndHostAndPortStr
		}
	}
	return uri.Validator()
}

func (uri *Uri) Validator() error {
	if reflect.DeepEqual(nil, uri) {
		return errors.New("uri caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(uri.schema)) == 0 {
		return errors.New("the schema field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(sip)$`).MatchString(uri.schema) {
		return errors.New("the value of the schema field must be sip")
	}
	if len(strings.TrimSpace(uri.user)) == 0 {
		return errors.New("the user field is not allowed to be empty")
	}
	if !regexp.MustCompile(`\d{20}`).MatchString(uri.user) {
		return errors.New("the user field must be 20 digits")
	}
	if len(strings.TrimSpace(uri.host)) == 0 {
		return errors.New("the host field is not allowed to be empty")
	}
	if ip := net.ParseIP(uri.host); ip != nil {
		if uri.port == 0 {
			return errors.New("the host field gives the ip address, the port must be given")
		}
	}
	return nil
}

func (uri *Uri) String() string {
	result := ""
	if len(strings.TrimSpace(uri.schema)) > 0 {
		result += fmt.Sprintf("%s:", strings.ToLower(uri.schema))
	}
	if len(strings.TrimSpace(uri.user)) > 0 {
		result += uri.user
	}
	if len(strings.TrimSpace(uri.host)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf("@%s", uri.host)
		} else {
			result += uri.host
		}
	}
	if uri.port > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(":%d", uri.port)
		} else {
			result += fmt.Sprintf("%d", uri.port)

		}
	}
	if !reflect.DeepEqual(nil, uri.extension) {
		extensions := ""
		for k, v := range uri.extension {
			if len(strings.TrimSpace(fmt.Sprintf("%v", v))) == 0 {
				extensions += fmt.Sprintf(";%s", k)
			} else {
				extensions += fmt.Sprintf(";%s=%v", k, v)
			}
		}
		if len(strings.TrimSpace(extensions)) > 0 {
			if len(result) == 0 {
				extensions = strings.TrimPrefix(extensions, ";")
			}
			result += extensions
		}
	}
	return result
}
