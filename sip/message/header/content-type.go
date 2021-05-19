package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ContentType struct {
	mediaType string // media-type
}

func (contentType *ContentType) SetMediaType(mediaType string) {
	contentType.mediaType = mediaType
}
func (contentType *ContentType) GetMediaType() string {
	return contentType.mediaType
}
func NewContentType(mediaType string) *ContentType {
	return &ContentType{
		mediaType: mediaType,
	}
}
func (contentType *ContentType) Raw() (string, error) {
	result := ""
	if err := contentType.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Content-Type: %s", contentType.mediaType)
	result += "\r\n"
	return result, nil
}
func (contentType *ContentType) Parse(raw string) error {
	if reflect.DeepEqual(nil, contentType) {
		return errors.New("content-type caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// content-type field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(content-type).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a content-type header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) > 0 {
		contentType.mediaType = raw
	}
	return contentType.Validator()
}
func (contentType *ContentType) Validator() error {
	if reflect.DeepEqual(nil, contentType) {
		return errors.New("content-type caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(contentType.mediaType)) == 0 {
		return errors.New("the media-type field is not allowed to be empty")
	}
	return nil
}
func (contentType *ContentType) String() string {
	result := ""
	if len(strings.TrimSpace(contentType.mediaType)) > 0 {
		result += contentType.mediaType
	}
	return result
}
