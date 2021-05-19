package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type CallID struct {
	id   string // ID
	host string // host
}

func (callId *CallID) SetId(id string) {
	callId.id = id
}
func (callId *CallID) GetId() string {
	return callId.id
}
func (callId *CallID) SetHost(host string) {
	callId.host = host
}
func (callId *CallID) GetHost() string {
	return callId.host
}
func NewCallID(id string, host string) *CallID {
	return &CallID{
		id:   id,
		host: host,
	}
}

func (callId *CallID) Raw() (string, error) {
	result := ""
	if err := callId.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(callId.host)) > 0 {
		result += fmt.Sprintf("Call-ID: %s@%s", callId.id, callId.host)
	} else {
		result += fmt.Sprintf("Call-ID: %s", callId.id)
	}
	result += "\r\n"
	return result, nil
}
func (callId *CallID) Parse(raw string) error {
	if reflect.DeepEqual(nil, callId) {
		return errors.New("call-id caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// call-id field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(call-id).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a call-id header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// host regexp
	hostRegexp := regexp.MustCompile(`@.*`)
	if hostRegexp.MatchString(raw) {
		callId.host = regexp.MustCompile(`@`).ReplaceAllString(hostRegexp.FindString(raw), "")
		raw = hostRegexp.ReplaceAllString(raw, "")
	}
	callId.id = raw
	return callId.Validator()
}
func (callId *CallID) Validator() error {
	if reflect.DeepEqual(nil, callId) {
		return errors.New("call-id caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(callId.id)) == 0 {
		return errors.New("the id field is not allowed to be empty")
	}
	return nil
}
func (callId *CallID) String() string {
	result := ""
	if len(strings.TrimSpace(callId.id)) > 0 {
		result += callId.id
	}
	if len(strings.TrimSpace(callId.host)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf("@%s", callId.host)
		} else {
			result += callId.host
		}
	}
	return result
}
