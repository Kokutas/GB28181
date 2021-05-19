package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type To struct {
	DisplayName string `json:"Display-Name"`
	Address     *Uri   `json:"SIP to Address"`
	Tag         string `json:"SIP to Tag"`
}

func (to *To) SetDisplayName(displayName string) {
	to.DisplayName = displayName
}
func (to *To) GetDisplayName() string {
	return to.DisplayName
}
func (to *To) SetAddress(address *Uri) {
	to.Address = address
}
func (to *To) GetAddress() *Uri {
	if to.Address != nil {
		return to.Address
	}
	return nil
}
func (to *To) SetTag(tag string) {
	to.Tag = tag
}
func (to *To) GetTag() string {
	return to.Tag
}

func NewTo(displayName string, address *Uri, tag string) *To {
	return &To{
		DisplayName: displayName,
		Address:     address,
		Tag:         tag,
	}
}
func (to *To) Raw() (string, error) {
	result := ""
	if err := to.Validator(); err != nil {
		return result, err
	}
	address, err := to.Address.Raw()
	if err != nil {
		return result, err
	}
	if len(strings.TrimSpace(to.DisplayName)) == 0 {
		result += fmt.Sprintf("To: <%s>", address)
	} else {
		result += fmt.Sprintf("To: \"%s\" %s", to.DisplayName, address)
	}
	if len(strings.TrimSpace(to.Tag)) > 0 {
		result += fmt.Sprintf(";tag=%s", to.Tag)
	}
	result += "\r\n"
	return result, nil
}
func (to *To) Parse(raw string) error {
	if reflect.DeepEqual(nil, to) {
		return errors.New("to caller is not allowed to be nil")
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
	// to field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(to).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a to header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// address and tag regexp
	addressAndTagRegexp := regexp.MustCompile(`(?i)(sip).*?:.*`)
	// tag regexp
	tagRegexp := regexp.MustCompile(`(?i)(tag)=.*`)
	// display-name
	displayNameStr := addressAndTagRegexp.ReplaceAllString(raw, "")
	displayNameStr = regexp.MustCompile(`<`).ReplaceAllString(displayNameStr, "")
	displayNameStr = regexp.MustCompile(`>`).ReplaceAllString(displayNameStr, "")
	displayNameStr = regexp.MustCompile(`"`).ReplaceAllString(displayNameStr, "")
	displayNameStr = strings.TrimLeft(displayNameStr, " ")
	displayNameStr = strings.TrimRight(displayNameStr, " ")
	displayNameStr = strings.TrimPrefix(displayNameStr, " ")
	displayNameStr = strings.TrimSuffix(displayNameStr, " ")
	if len(strings.TrimSpace(displayNameStr)) > 0 {
		to.DisplayName = displayNameStr
	}
	if addressAndTagRegexp.MatchString(raw) {
		addressAndTag := addressAndTagRegexp.FindString(raw)
		if tagRegexp.MatchString(addressAndTag) {
			to.Tag = regexp.MustCompile(`(?i)(tag)=`).ReplaceAllString(tagRegexp.FindString(addressAndTag), "")
			addressAndTag = tagRegexp.ReplaceAllString(addressAndTag, "")
		}
		addressAndTag = regexp.MustCompile(`<`).ReplaceAllString(addressAndTag, "")
		addressAndTag = regexp.MustCompile(`>`).ReplaceAllString(addressAndTag, "")
		addressAndTag = strings.TrimLeft(addressAndTag, ";")
		addressAndTag = strings.TrimRight(addressAndTag, ";")
		addressAndTag = strings.TrimPrefix(addressAndTag, ";")
		addressAndTag = strings.TrimSuffix(addressAndTag, ";")
		addressAndTag = strings.TrimLeft(addressAndTag, " ")
		addressAndTag = strings.TrimRight(addressAndTag, " ")
		addressAndTag = strings.TrimPrefix(addressAndTag, " ")
		addressAndTag = strings.TrimSuffix(addressAndTag, " ")
		to.Address = new(Uri)
		if err := to.Address.Parse(addressAndTag); err != nil {
			return err
		}
	}

	return to.Validator()
}
func (to *To) Validator() error {
	if reflect.DeepEqual(nil, to) {
		return errors.New("to caller is not allowed to be nil")
	}
	if err := to.Address.Validator(); err != nil {
		return fmt.Errorf("to address validator error : %s", err.Error())
	}
	return nil
}
