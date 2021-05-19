package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type From struct {
	displayName string // display-name
	address     *Uri   // SIP from Address"`
	tag         string // SIP from Tag"`
}

func (from *From) SetDisplayName(displayName string) {
	from.displayName = displayName
}
func (from *From) GetDisplayName() string {
	return from.displayName
}
func (from *From) SetAddress(address *Uri) {
	from.address = address
}
func (from *From) GetAddress() *Uri {
	return from.address
}
func (from *From) SetTag(tag string) {
	from.tag = tag
}
func (from *From) GetTag() string {
	return from.tag
}

func NewFrom(displayName string, address *Uri, tag string) *From {
	return &From{
		displayName: displayName,
		address:     address,
		tag:         tag,
	}
}
func (from *From) Raw() (string, error) {
	result := ""
	if err := from.Validator(); err != nil {
		return result, err
	}
	address, err := from.address.Raw()
	if err != nil {
		return result, err
	}
	if len(strings.TrimSpace(from.displayName)) == 0 {
		result += fmt.Sprintf("From: <%s>", address)
	} else {
		result += fmt.Sprintf("From: \"%s\" %s", from.displayName, address)
	}
	if len(strings.TrimSpace(from.tag)) > 0 {
		result += fmt.Sprintf(";tag=%s", from.tag)
	}
	result += "\r\n"
	return result, nil
}
func (from *From) Parse(raw string) error {
	if reflect.DeepEqual(nil, from) {
		return errors.New("from caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// from field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(from).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a from header field")
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
		from.displayName = displayNameStr
	}
	if addressAndTagRegexp.MatchString(raw) {
		addressAndTag := addressAndTagRegexp.FindString(raw)
		if tagRegexp.MatchString(addressAndTag) {
			from.tag = regexp.MustCompile(`(?i)(tag)=`).ReplaceAllString(tagRegexp.FindString(addressAndTag), "")
			addressAndTag = tagRegexp.ReplaceAllString(addressAndTag, "")
		}
		addressAndTag = regexp.MustCompile(`<`).ReplaceAllString(addressAndTag, "")
		addressAndTag = regexp.MustCompile(`>`).ReplaceAllString(addressAndTag, "")
		addressAndTag = strings.TrimPrefix(addressAndTag, ";")
		addressAndTag = strings.TrimSuffix(addressAndTag, ";")
		addressAndTag = strings.TrimPrefix(addressAndTag, " ")
		addressAndTag = strings.TrimSuffix(addressAndTag, " ")
		from.address = new(Uri)
		if err := from.address.Parse(addressAndTag); err != nil {
			return err
		}
	}

	return from.Validator()
}
func (from *From) Validator() error {
	if reflect.DeepEqual(nil, from) {
		return errors.New("from caller is not allowed to be nil")
	}
	if err := from.address.Validator(); err != nil {
		return fmt.Errorf("from address validator error : %s", err.Error())
	}
	return nil
}
func (from *From) String() string {
	result := ""
	if len(strings.TrimSpace(from.displayName)) > 0 {
		result += fmt.Sprintf(" \"%s\"", from.displayName)
		if !reflect.DeepEqual(nil, from.address) {
			if len(result) > 0 {
				result += fmt.Sprintf(" %s", from.address.String())
			}
		}
	} else {
		if !reflect.DeepEqual(nil, from.address) {
			if len(result) > 0 {
				result += fmt.Sprintf("<%s>", from.address.String())
			}
		}
	}

	if len(strings.TrimSpace(from.tag)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(";tag=%s", from.tag)
		} else {
			result += fmt.Sprintf("tag=%s", from.tag)
		}
	}
	return result
}
