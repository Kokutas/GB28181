package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type CallID struct {
	ID   string `json:"ID"`
	Host string `json:"Host"`
}

func (callId *CallID) SetID(id string) {
	callId.ID = id
}
func (callId *CallID) GetID() string {
	return callId.ID
}
func (callId *CallID) SetHost(host string) {
	callId.Host = host
}
func (callId *CallID) GetHost() string {
	return callId.Host
}
func NewCallID(id string, host string) *CallID {
	return &CallID{
		ID:   id,
		Host: host,
	}
}

func (callId *CallID) Raw() (string, error) {
	result := ""
	if err := callId.Validator(); err != nil {
		return result, err
	}
	if len(strings.TrimSpace(callId.Host)) > 0 {
		result += fmt.Sprintf("Call-ID: %s@%s", callId.ID, callId.Host)
	} else {
		result += fmt.Sprintf("Call-ID: %s", callId.ID)
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
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// call-id field regexp
	callIdFieldRegexp := regexp.MustCompile(`(?i)(call-id).*?:`)
	if !callIdFieldRegexp.MatchString(raw) {
		return errors.New("raw is not a call-id header field")
	}
	raw = callIdFieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// host regexp
	hostRegexp := regexp.MustCompile(`@.*`)
	if hostRegexp.MatchString(raw) {
		callId.Host = regexp.MustCompile(`@`).ReplaceAllString(hostRegexp.FindString(raw), "")
		raw = hostRegexp.ReplaceAllString(raw, "")
	}
	callId.ID = raw
	return callId.Validator()
}
func (callId *CallID) Validator() error {
	if reflect.DeepEqual(nil, callId) {
		return errors.New("call-id caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(callId.ID)) == 0 {
		return errors.New("the id field is not allowed to be empty")
	}
	return nil
}
