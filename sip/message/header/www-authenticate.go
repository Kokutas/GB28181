package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type WWWAuthenticate struct {
	AuthSchema string `json:"AuthSchema"` // basic / digest
	Realm      string `json:"Realm"`
	Nonce      string `json:"Nonce"`
	Algorithm  string `json:"Algorithm"`
}

func (wwwAuthenticate *WWWAuthenticate) SetAuthSchema(authSchema string) {
	wwwAuthenticate.AuthSchema = authSchema
}
func (wwwAuthenticate *WWWAuthenticate) GetAuthSchema() string {
	return wwwAuthenticate.AuthSchema
}
func (wwwAuthenticate *WWWAuthenticate) SetRealm(realm string) {
	wwwAuthenticate.Realm = realm
}
func (wwwAuthenticate *WWWAuthenticate) GetRealm() string {
	return wwwAuthenticate.Realm
}
func (wwwAuthenticate *WWWAuthenticate) SetNonce(nonce string) {
	wwwAuthenticate.Nonce = nonce
}
func (wwwAuthenticate *WWWAuthenticate) GetNonce() string {
	return wwwAuthenticate.Nonce
}

func (wwwAuthenticate *WWWAuthenticate) SetAlgorithm(algorithm string) {
	wwwAuthenticate.Algorithm = algorithm
}
func (wwwAuthenticate *WWWAuthenticate) GetAlgorithm() string {
	return wwwAuthenticate.Algorithm
}
func NewWWWAuthenticate(authSchema string, realm string, nonce string, algorithm string) *WWWAuthenticate {
	return &WWWAuthenticate{
		AuthSchema: authSchema,
		Realm:      realm,
		Nonce:      nonce,
		Algorithm:  algorithm,
	}
}
func (wwwAuthenticate *WWWAuthenticate) Raw() (string, error) {
	result := ""
	if err := wwwAuthenticate.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("WWW-Authenticate: %s realm=\"%s\",nonce=\"%s\"", strings.Title(wwwAuthenticate.AuthSchema), wwwAuthenticate.Realm, wwwAuthenticate.Nonce)
	if len(strings.TrimSpace(wwwAuthenticate.Algorithm)) > 0 {
		result += fmt.Sprintf(",algorithm=%s", strings.ToUpper(wwwAuthenticate.Algorithm))
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
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
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
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		wwwAuthenticate.AuthSchema = authSchemaRegexp.FindString(raw)
	}
	raw = authSchemaRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// realm regexp
	realmRegexp := regexp.MustCompile(`(?i)(realm).*?=`)
	// nonce regexp
	nonceRegexp := regexp.MustCompile(`(?i)(nonce).*?=`)
	// algorithm regexp
	algorithmRegexp := regexp.MustCompile(`(?i)(algorithm).*?=`)
	raw = strings.TrimLeft(raw, ",")
	raw = strings.TrimRight(raw, ",")
	raw = strings.TrimPrefix(raw, ",")
	raw = strings.TrimSuffix(raw, ",")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case realmRegexp.MatchString(raws):
			wwwAuthenticate.Realm = realmRegexp.ReplaceAllString(raws, "")
		case nonceRegexp.MatchString(raws):
			wwwAuthenticate.Nonce = nonceRegexp.ReplaceAllString(raws, "")
		case algorithmRegexp.MatchString(raws):
			wwwAuthenticate.Algorithm = algorithmRegexp.ReplaceAllString(raws, "")
		}
	}
	return wwwAuthenticate.Validator()

}
func (wwwAuthenticate *WWWAuthenticate) Validator() error {
	if reflect.DeepEqual(nil, wwwAuthenticate) {
		return errors.New("www-authenticate caller is not allowed to be nil")
	}
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if !authSchemaRegexp.MatchString(wwwAuthenticate.AuthSchema) {
		return errors.New("the value of the authschema field must be Digest")
	}
	if len(strings.TrimSpace(wwwAuthenticate.Realm)) == 0 {
		return errors.New("the realm field is not allowed to be empty")
	}
	if len(strings.TrimSpace(wwwAuthenticate.Nonce)) == 0 {
		return errors.New("the nonce field is not allowed to be empty")
	}
	if len(strings.TrimSpace(wwwAuthenticate.Algorithm)) > 0 {
		if !regexp.MustCompile(`(?i)(md5)`).MatchString(wwwAuthenticate.Algorithm) {
			return errors.New("the value of the algorithm field must be MD5")
		}
	}
	return nil
}
