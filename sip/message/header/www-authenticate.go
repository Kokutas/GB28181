package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type WWWAuthenticate struct {
	authSchema string // authSchema // basic / digest
	realm      string // realm
	nonce      string // nonce
	algorithm  string // algorithm
}

func (wwwAuthenticate *WWWAuthenticate) SetAuthSchema(authSchema string) {
	wwwAuthenticate.authSchema = authSchema
}
func (wwwAuthenticate *WWWAuthenticate) GetAuthSchema() string {
	return wwwAuthenticate.authSchema
}
func (wwwAuthenticate *WWWAuthenticate) SetRealm(realm string) {
	wwwAuthenticate.realm = realm
}
func (wwwAuthenticate *WWWAuthenticate) GetRealm() string {
	return wwwAuthenticate.realm
}
func (wwwAuthenticate *WWWAuthenticate) SetNonce(nonce string) {
	wwwAuthenticate.nonce = nonce
}
func (wwwAuthenticate *WWWAuthenticate) GetNonce() string {
	return wwwAuthenticate.nonce
}

func (wwwAuthenticate *WWWAuthenticate) SetAlgorithm(algorithm string) {
	wwwAuthenticate.algorithm = algorithm
}
func (wwwAuthenticate *WWWAuthenticate) GetAlgorithm() string {
	return wwwAuthenticate.algorithm
}
func NewWWWAuthenticate(authSchema string, realm string, nonce string, algorithm string) *WWWAuthenticate {
	return &WWWAuthenticate{
		authSchema: authSchema,
		realm:      realm,
		nonce:      nonce,
		algorithm:  algorithm,
	}
}
func (wwwAuthenticate *WWWAuthenticate) Raw() (string, error) {
	result := ""
	if err := wwwAuthenticate.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("WWW-Authenticate: %s realm=\"%s\",nonce=\"%s\"", strings.Title(wwwAuthenticate.authSchema), wwwAuthenticate.realm, wwwAuthenticate.nonce)
	if len(strings.TrimSpace(wwwAuthenticate.algorithm)) > 0 {
		result += fmt.Sprintf(",algorithm=%s", strings.ToUpper(wwwAuthenticate.algorithm))
	}
	result += "\r\n"
	return result, nil
}
func (wwwAuthenticate *WWWAuthenticate) Parse(raw string) error {
	if reflect.DeepEqual(nil, wwwAuthenticate) {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// www-authenticate field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(www-authenticate).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a www-authenticate header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		wwwAuthenticate.authSchema = authSchemaRegexp.FindString(raw)
	}
	raw = authSchemaRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// realm regexp
	realmRegexp := regexp.MustCompile(`(?i)(realm).*?=`)
	// nonce regexp
	nonceRegexp := regexp.MustCompile(`(?i)(nonce).*?=`)
	// algorithm regexp
	algorithmRegexp := regexp.MustCompile(`(?i)(algorithm).*?=`)
	raw = strings.TrimPrefix(raw, ",")
	raw = strings.TrimSuffix(raw, ",")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case realmRegexp.MatchString(raws):
			wwwAuthenticate.realm = realmRegexp.ReplaceAllString(raws, "")
		case nonceRegexp.MatchString(raws):
			wwwAuthenticate.nonce = nonceRegexp.ReplaceAllString(raws, "")
		case algorithmRegexp.MatchString(raws):
			wwwAuthenticate.algorithm = algorithmRegexp.ReplaceAllString(raws, "")
		}
	}
	return wwwAuthenticate.Validator()

}
func (wwwAuthenticate *WWWAuthenticate) Validator() error {
	if reflect.DeepEqual(nil, wwwAuthenticate) {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if !authSchemaRegexp.MatchString(wwwAuthenticate.authSchema) {
		return errors.New("the value of the authschema field must be Digest")
	}
	if len(strings.TrimSpace(wwwAuthenticate.realm)) == 0 {
		return errors.New("the realm field is not allowed to be empty")
	}
	if len(strings.TrimSpace(wwwAuthenticate.nonce)) == 0 {
		return errors.New("the nonce field is not allowed to be empty")
	}
	if len(strings.TrimSpace(wwwAuthenticate.algorithm)) > 0 {
		if !regexp.MustCompile(`(?i)(md5)`).MatchString(wwwAuthenticate.algorithm) {
			return errors.New("the value of the algorithm field must be MD5")
		}
	}
	return nil
}

func (wwwAuthenticate *WWWAuthenticate) String() string {
	result := ""
	if len(strings.TrimSpace(wwwAuthenticate.authSchema)) > 0 {
		result += fmt.Sprintf("%s", strings.Title(wwwAuthenticate.authSchema))
	}
	if len(strings.TrimSpace(wwwAuthenticate.realm)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(",realm=\"%s\"", wwwAuthenticate.realm)
		} else {
			result += fmt.Sprintf("realm=\"%s\"", wwwAuthenticate.realm)
		}
	}
	if len(strings.TrimSpace(wwwAuthenticate.nonce)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(",nonce=\"%s\"", wwwAuthenticate.nonce)
		} else {
			result += fmt.Sprintf("nonce=\"%s\"", wwwAuthenticate.nonce)
		}
	}
	if len(strings.TrimSpace(wwwAuthenticate.algorithm)) > 0 {
		if len(result) > 0 {
			result += fmt.Sprintf(",algorithm=%s", strings.ToUpper(wwwAuthenticate.algorithm))
		} else {
			result += fmt.Sprintf("algorithm=%s", strings.ToUpper(wwwAuthenticate.algorithm))
		}
	}
	return result
}
