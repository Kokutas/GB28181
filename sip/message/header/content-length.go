package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ContentLength struct {
	length uint // body-content-length
}

func (contentLength *ContentLength) SetLength(length uint) {
	contentLength.length = length
}

func (contentLength *ContentLength) GetLength() uint {
	return contentLength.length
}

func NewContentLength(length uint) *ContentLength {
	return &ContentLength{
		length: length,
	}
}
func (contentLength *ContentLength) Raw() (string, error) {
	result := ""
	if err := contentLength.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Content-Length: %d", contentLength.length)
	result += "\r\n"
	return result, nil
}
func (contentLength *ContentLength) Parse(raw string) error {
	if reflect.DeepEqual(nil, contentLength) {
		return errors.New("content-length caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// content-length field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(content-length).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a content-length header field")
	}
	// length regexp
	lengthRegexp := regexp.MustCompile(`\d+`)
	if lengthRegexp.MatchString(raw) {
		lengthStr := lengthRegexp.FindString(raw)
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return err
		}
		contentLength.length = uint(length)
	}

	return contentLength.Validator()
}
func (contentLength *ContentLength) Validator() error {
	if reflect.DeepEqual(nil, contentLength) {
		return errors.New("content-length caller is not allowed to be nil")
	}
	return nil
}

func (contentLength *ContentLength) String() string {
	result := ""
	result += fmt.Sprintf("%d", contentLength.length)
	return result
}
