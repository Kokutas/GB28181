package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Contact struct {
	DisplayName string `json:"Display-Name"`
	*Uri        `json:"Uri"`
	Extension   map[string]interface{} `json:"Extension"`
}

func (contact *Contact) SetDisplayName(displayName string) {
	contact.DisplayName = displayName
}
func (contact *Contact) GetDisplayName() string {
	return contact.DisplayName
}
func (contact *Contact) SetUri(uri *Uri) {
	contact.Uri = uri
}
func (contact *Contact) GetUri() *Uri {
	if contact.Uri != nil {
		return contact.Uri
	}
	return nil
}
func (contact *Contact) SetExtension(extensions map[string]interface{}) {
	contact.Extension = extensions
}
func (contact *Contact) GetExtension() map[string]interface{} {
	if contact.Extension != nil {
		return contact.Extension
	}
	return nil
}

func NewContact(displayName string, uri *Uri, extension map[string]interface{}) *Contact {
	return &Contact{
		DisplayName: displayName,
		Uri:         uri,
		Extension:   extension,
	}
}
func (contact *Contact) Raw() (string, error) {
	result := ""
	if err := contact.Validator(); err != nil {
		return result, err
	}
	uri, err := contact.Uri.Raw()
	if err != nil {
		return result, err
	}
	if len(strings.TrimSpace(contact.DisplayName)) == 0 {
		result += fmt.Sprintf("Contact: <%s>", uri)
	} else {
		result += fmt.Sprintf("Contact: \"%s\" %s", contact.DisplayName, uri)
	}
	if contact.Extension != nil {
		extensions := ""
		for k, v := range contact.Extension {
			if len(strings.TrimSpace(fmt.Sprintf("%v", v))) == 0 {
				extensions += fmt.Sprintf(";%s", k)
			} else {
				extensions += fmt.Sprintf(";%s=%v", k, v)
			}
		}
		if len(strings.TrimSpace(extensions)) > 0 {
			result += extensions
		}
	}
	result += "\r\n"
	return result, nil
}
func (contact *Contact) Parse(raw string) error {
	if reflect.DeepEqual(nil, contact) {
		return errors.New("contact caller is not allowed to be nil")
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
	// contact field regexp
	contactFieldRegexp := regexp.MustCompile(`(?i)(contact).*?:`)
	if !contactFieldRegexp.MatchString(raw) {
		return errors.New("raw is not a contact header field")
	}
	raw = contactFieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// extensions regexp
	extensionsRegexp := regexp.MustCompile(`>;.*`)
	if extensionsRegexp.MatchString(raw) {
		m := make(map[string]interface{})
		extension := extensionsRegexp.FindString(raw)
		extension = regexp.MustCompile(`>`).ReplaceAllString(extension, "")
		extension = strings.TrimLeft(extension, ";")
		extension = strings.TrimRight(extension, ";")
		extension = strings.TrimPrefix(extension, ";")
		extension = strings.TrimSuffix(extension, ";")
		extensions := strings.Split(extension, ";")
		for _, v := range extensions {
			if strings.Contains(v, "=") {
				vs := strings.Split(v, "=")
				if len(vs) > 1 {
					m[vs[0]] = vs[1]
				} else {
					m[vs[0]] = ""
				}
			} else {
				m[v] = ""
			}
		}
		if len(m) > 0 {
			contact.Extension = m
		}
	}
	raw = extensionsRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// uri regexp
	uriRegex := regexp.MustCompile(`(?i)(sip).*?:.*`)
	if uriRegex.MatchString(raw) {
		uris := uriRegex.FindString(raw)
		uris = regexp.MustCompile(`>`).ReplaceAllString(uris, "")
		uris = regexp.MustCompile(`<`).ReplaceAllString(uris, "")
		contact.Uri = new(Uri)
		if err := contact.Uri.Parse(uris); err != nil {
			return err
		}
	}

	raw = uriRegex.ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`>`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`<`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`"`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`"`).ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// display name regexp
	if len(strings.TrimSpace(raw)) > 0 {
		contact.DisplayName = raw
	}

	return contact.Validator()
}
func (contact *Contact) Validator() error {
	if reflect.DeepEqual(nil, contact) {
		return errors.New("contact caller is not allowed to be nil")
	}
	if err := contact.Uri.Validator(); err != nil {
		return fmt.Errorf("contact uri validator error : %s", err.Error())
	}
	return nil
}
