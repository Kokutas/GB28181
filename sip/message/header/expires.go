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
	seconds uint // seconds
}

func (expires *Expires) SetSeconds(seconds uint) {
	expires.seconds = seconds
}
func (expires *Expires) GetSeconds() uint {
	return expires.seconds
}
func NewExpires(seconds uint) *Expires {
	return &Expires{
		seconds: seconds,
	}
}

func (expires *Expires) Raw() (string, error) {
	result := ""
	if err := expires.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Expires: %d", expires.seconds)
	result += "\r\n"
	return result, nil
}

func (expires *Expires) Parse(raw string) error {
	if reflect.DeepEqual(nil, expires) {
		return errors.New("expires caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// expires field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(expires).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a expires header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
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
		expires.seconds = uint(second)
	}
	return nil
}
func (expires *Expires) Validator() error {
	if reflect.DeepEqual(nil, expires) {
		return errors.New("expires caller is not allowed to be nil")
	}
	return nil
}
func (expires *Expires) String() string {
	result := ""
	result += fmt.Sprintf("%d", expires.seconds)
	return result
}
