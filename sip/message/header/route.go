package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	DisplayName string `json:"Display-Name"`
	Uris        []*Uri `json:"Route URI"`
}

func (route *Route) SetDisplayName(displayName string) {
	route.DisplayName = displayName
}
func (route *Route) GetDisplayName() string {
	return route.DisplayName
}

func (route *Route) SetUris(uris ...*Uri) {
	route.Uris = append(route.Uris, uris...)
}
func (route *Route) GetUris() []*Uri {
	if len(route.Uris) > 0 {
		return route.Uris
	}
	return nil
}
func NewRoute(displayName string, uris ...*Uri) *Route {
	return &Route{
		DisplayName: displayName,
		Uris:        uris,
	}
}

func (route *Route) Raw() (string, error) {
	result := ""
	if err := route.Validator(); err != nil {
		return result, err
	}
	result += "Route:"
	addQuoteTag := true
	if len(strings.TrimSpace(route.DisplayName)) > 0 {
		addQuoteTag = false
		result += fmt.Sprintf(" \"%s\"", route.DisplayName)
	}
	for _, uri := range route.Uris {
		uriStr, err := uri.Raw()
		if err != nil {
			return "", err
		}
		if !addQuoteTag {
			result += fmt.Sprintf(" %s,", uriStr)
		} else {
			result += fmt.Sprintf(" <%s>,", uriStr)
		}
	}
	result = strings.TrimSuffix(result, ",")
	result += "\r\n"
	return result, nil
}
func (route *Route) Parse(raw string) error {
	if reflect.DeepEqual(nil, route) {
		return errors.New("route caller is not allowed to be nil")
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
	// route field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(route).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a route header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// address and tag regexp
	addressAndTagRegexp := regexp.MustCompile(`(?i)(sip).*?:.*`)
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
		route.DisplayName = displayNameStr
		raw = regexp.MustCompile(`.*`+displayNameStr).ReplaceAllString(raw, "")
	}
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	raw = strings.TrimLeft(raw, ",")
	raw = strings.TrimRight(raw, ",")
	raw = strings.TrimPrefix(raw, ",")
	raw = strings.TrimSuffix(raw, ",")

	if strings.Contains(raw, ",") {
		rawSlice := strings.Split(raw, ",")
		route.Uris = make([]*Uri, 0, len(rawSlice))
		for _, raws := range rawSlice {
			raws = regexp.MustCompile(`<`).ReplaceAllString(raws, "")
			raws = regexp.MustCompile(`>`).ReplaceAllString(raws, "")
			uri := new(Uri)
			if err := uri.Parse(raws); err != nil {
				return err
			}
			route.Uris = append(route.Uris, uri)
		}

	} else {
		raw = regexp.MustCompile(`<`).ReplaceAllString(raw, "")
		raw = regexp.MustCompile(`>`).ReplaceAllString(raw, "")
		route.Uris = make([]*Uri, 0, 1)
		uri := new(Uri)
		if err := uri.Parse(raw); err != nil {
			return err
		}
		route.Uris = append(route.Uris, uri)
	}

	return route.Validator()
}
func (route *Route) Validator() error {
	if reflect.DeepEqual(nil, route) {
		return errors.New("route caller is not allowed to be nil")
	}
	if reflect.DeepEqual(nil, route.Uris) {
		return errors.New("the uris field is not allowed to be nil")
	}
	if len(route.Uris) == 0 {
		return errors.New("the uris field must has one uri")
	}
	return nil
}
