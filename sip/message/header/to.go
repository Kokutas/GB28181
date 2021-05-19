package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type To struct {
	displayName string // display-name
	address     *Uri   // SIP to Address
	tag         string // SIP to Tag
}

func (to *To) SetDisplayName(displayName string) {
	to.displayName = displayName
}
func (to *To) GetDisplayName() string {
	return to.displayName
}
func (to *To) SetAddress(address *Uri) {
	to.address = address
}
func (to *To) GetAddress() *Uri {
	return to.address
}
func (to *To) SetTag(tag string) {
	to.tag = tag
}
func (to *To) GetTag() string {
	return to.tag
}

func NewTo(displayName string, address *Uri, tag string) *To {
	return &To{
		displayName: displayName,
		address:     address,
		tag:         tag,
	}
}
func (to *To) Raw() (string, error) {
	result := ""
	if err := to.Validator(); err != nil {
		return result, err
	}
	address, err := to.address.Raw()
	if err != nil {
		return result, err
	}
	if len(strings.TrimSpace(to.displayName)) == 0 {
		result += fmt.Sprintf("To: <%s>", address)
	} else {
		result += fmt.Sprintf("To: \"%s\" %s", to.displayName, address)
	}
	if len(strings.TrimSpace(to.tag)) > 0 {
		result += fmt.Sprintf(";tag=%s", to.tag)
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
	displayNameStr = strings.TrimPrefix(displayNameStr, " ")
	displayNameStr = strings.TrimSuffix(displayNameStr, " ")
	if len(strings.TrimSpace(displayNameStr)) > 0 {
		to.displayName = displayNameStr
	}
	if addressAndTagRegexp.MatchString(raw) {
		addressAndTag := addressAndTagRegexp.FindString(raw)
		if tagRegexp.MatchString(addressAndTag) {
			to.tag = regexp.MustCompile(`(?i)(tag)=`).ReplaceAllString(tagRegexp.FindString(addressAndTag), "")
			addressAndTag = tagRegexp.ReplaceAllString(addressAndTag, "")
		}
		addressAndTag = regexp.MustCompile(`<`).ReplaceAllString(addressAndTag, "")
		addressAndTag = regexp.MustCompile(`>`).ReplaceAllString(addressAndTag, "")
		addressAndTag = strings.TrimPrefix(addressAndTag, ";")
		addressAndTag = strings.TrimSuffix(addressAndTag, ";")
		addressAndTag = strings.TrimPrefix(addressAndTag, " ")
		addressAndTag = strings.TrimSuffix(addressAndTag, " ")
		to.address = new(Uri)
		if err := to.address.Parse(addressAndTag); err != nil {
			return err
		}
	}

	return to.Validator()
}
func (to *To) Validator() error {
	if reflect.DeepEqual(nil, to) {
		return errors.New("to caller is not allowed to be nil")
	}
	if err := to.address.Validator(); err != nil {
		return fmt.Errorf("to address validator error : %s", err.Error())
	}
	return nil
}
func (to *To) String() string {
	result := ""
	if len(strings.TrimSpace(to.displayName)) > 0 {
		result += fmt.Sprintf(" \"%s\"", to.displayName)
		if !reflect.DeepEqual(nil, to.address) {
			if len(result) > 0 {
				result += fmt.Sprintf(" %s", to.address.String())
			}
		}
	} else {
		if !reflect.DeepEqual(nil, to.address) {
			if len(result) > 0 {
				result += fmt.Sprintf("<%s>", to.address.String())
			}
		}
	}

	if len(strings.TrimSpace(to.tag)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(";tag=%s", to.tag)
		} else {
			result += fmt.Sprintf("tag=%s", to.tag)
		}
	}
	return result
}
