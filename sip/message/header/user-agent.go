package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type UserAgent struct {
	server string // server
}

func (userAgent *UserAgent) SetServer(server string) {
	userAgent.server = server
}
func (userAgent *UserAgent) GetServer() string {
	return userAgent.server
}
func NewUserAgent(server string) *UserAgent {
	return &UserAgent{
		server: server,
	}
}

func (userAgent *UserAgent) Raw() (string, error) {
	result := ""
	if err := userAgent.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("User-Agent: %s", userAgent.server)
	result += "\r\n"
	return result, nil
}
func (userAgent *UserAgent) Parse(raw string) error {
	if reflect.DeepEqual(nil, userAgent) {
		return errors.New("user-agent caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// user-agent field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(user-agent).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a user-agent header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	if len(strings.TrimSpace(raw)) > 0 {
		userAgent.server = raw
	}
	return userAgent.Validator()
}
func (userAgent *UserAgent) Validator() error {
	if reflect.DeepEqual(nil, userAgent) {
		return errors.New("user-agent caller is not allowed to be nil")
	}
	return nil
}
func (userAgent *UserAgent) String() string {
	result := ""
	if len(strings.TrimSpace(userAgent.server)) > 0 {
		result += userAgent.server
	}
	return result
}
