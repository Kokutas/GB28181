package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Expires struct {
	Seconds uint `json:"Seconds"`
}

func (expires *Expires) SetSeconds(seconds uint) {
	expires.Seconds = seconds
}
func (expires *Expires) GetSeconds() uint {
	return expires.Seconds
}
func NewExpires(seconds uint) *Expires {
	return &Expires{
		Seconds: seconds,
	}
}

func (expires *Expires) Raw() (string, error) {
	result := ""
	if err := expires.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Expires: %d", expires.Seconds)
	result += "\r\n"
	return result, nil
}

func (expires *Expires) Parse(raw string) error {
	if reflect.DeepEqual(nil, expires) {
		return errors.New("expires caller is not allowed to be nil")
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
	// expires field regexp
	expiresFieldRegexp := regexp.MustCompile(`(?i)(expires).*?:`)
	if !expiresFieldRegexp.MatchString(raw) {
		return errors.New("raw is not a expires header field")
	}
	raw = expiresFieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// seconds regexp
	secondsRegexp := regexp.MustCompile(`\d+`)
	if secondsRegexp.MatchString(raw) {
		secondStr := secondsRegexp.FindString(raw)
		second, err := strconv.Atoi(secondStr)
		if err != nil {
			return err
		}
		expires.Seconds = uint(second)
	}
	return nil
}
func (expires *Expires) Validator() error {
	if reflect.DeepEqual(nil, expires) {
		return errors.New("expires caller is not allowed to be nil")
	}
	return nil
}
