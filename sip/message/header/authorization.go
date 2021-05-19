package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Authorization struct {
	AuthSchema string `json:"AuthSchema"` // basic / digest
	UserName   string `json:"UserName"`
	Realm      string `json:"Realm"`
	Nonce      string `json:"Nonce"`
	*Uri       `json:"Uri"`
	Response   string `json:"Response"`
	Algorithm  string `json:"Algorithm"`
}

func (authorization *Authorization) SetAuthSchema(authSchema string) {
	authorization.AuthSchema = authSchema
}
func (authorization *Authorization) GetAuthSchema() string {
	return authorization.AuthSchema
}
func (authorization *Authorization) SetUserName(username string) {
	authorization.UserName = username
}
func (authorization *Authorization) GetUserName() string {
	return authorization.UserName
}
func (authorization *Authorization) SetRealm(realm string) {
	authorization.Realm = realm
}
func (authorization *Authorization) GetRealm() string {
	return authorization.Realm
}
func (authorization *Authorization) SetNonce(nonce string) {
	authorization.Nonce = nonce
}
func (authorization *Authorization) GetNonce() string {
	return authorization.Nonce
}

func (authorization *Authorization) SetUri(uri *Uri) {
	authorization.Uri = uri
}
func (authorization *Authorization) GetUri() *Uri {
	if authorization.Uri != nil {
		return authorization.Uri
	}
	return nil
}
func (authorization *Authorization) SetResponse(response string) {
	authorization.Response = response
}
func (authorization *Authorization) GetResponse() string {
	return authorization.Response
}
func (authorization *Authorization) SetAlgorithm(algorithm string) {
	authorization.Algorithm = algorithm
}
func (authorization *Authorization) GetAlgorithm() string {
	return authorization.Algorithm
}
func NewAuthorization(authSchema string, username string, realm string, nonce string, uri *Uri, response string, algorithm string) *Authorization {
	return &Authorization{
		AuthSchema: authSchema,
		UserName:   username,
		Realm:      realm,
		Nonce:      nonce,
		Uri:        uri,
		Response:   response,
		Algorithm:  algorithm,
	}
}

func (authorization *Authorization) Raw() (string, error) {
	result := ""
	if err := authorization.Validator(); err != nil {
		return result, err
	}
	uriStr, err := authorization.Uri.Raw()
	if err != nil {
		return result, err
	}
	result += fmt.Sprintf("Authorization: %s username=\"%s\",realm=\"%s\",nonce=\"%s\",uri=\"%s\",response=\"%s\",algorithm=%s",
		strings.Title(authorization.AuthSchema), authorization.UserName, authorization.Realm, authorization.Nonce, uriStr, authorization.Response, strings.ToUpper(authorization.Algorithm))
	result += "\r\n"
	return result, nil
}
func (authorization *Authorization) Parse(raw string) error {
	if reflect.DeepEqual(nil, authorization) {
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
	// authorization field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(Authorization).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a authorization header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")

	// auth schema regexp
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if authSchemaRegexp.MatchString(raw) {
		authorization.AuthSchema = authSchemaRegexp.FindString(raw)
	}
	raw = authSchemaRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimLeft(raw, " ")
	raw = strings.TrimRight(raw, " ")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// username regexp
	usernameRegexp := regexp.MustCompile(`(?i)(username).*?=`)
	// realm regexp
	realmRegexp := regexp.MustCompile(`(?i)(realm).*?=`)
	// nonce regexp
	nonceRegexp := regexp.MustCompile(`(?i)(nonce).*?=`)
	// uri regexp
	uriRegexp := regexp.MustCompile(`(?i)(uri).*?=`)
	// response regexp
	responseRegexp := regexp.MustCompile(`(?i)(response).*?=`)
	// algorithm regexp
	algorithmRegexp := regexp.MustCompile(`(?i)(algorithm).*?=`)
	raw = strings.TrimLeft(raw, ",")
	raw = strings.TrimRight(raw, ",")
	raw = strings.TrimPrefix(raw, ",")
	raw = strings.TrimSuffix(raw, ",")
	rawSlice := strings.Split(raw, ",")
	for _, raws := range rawSlice {
		switch {
		case usernameRegexp.MatchString(raws):
			authorization.UserName = usernameRegexp.ReplaceAllString(raws, "")
		case realmRegexp.MatchString(raws):
			authorization.Realm = realmRegexp.ReplaceAllString(raws, "")
		case nonceRegexp.MatchString(raws):
			authorization.Nonce = nonceRegexp.ReplaceAllString(raws, "")
		case uriRegexp.MatchString(raws):
			authorization.Uri = new(Uri)
			uriStr := uriRegexp.ReplaceAllString(raws, "")
			uriStr = regexp.MustCompile(`"`).ReplaceAllString(uriStr, "")
			uriStr = regexp.MustCompile(`<`).ReplaceAllString(uriStr, "")
			uriStr = regexp.MustCompile(`>`).ReplaceAllString(uriStr, "")
			if err := authorization.Uri.Parse(uriStr); err != nil {
				return err
			}
		case responseRegexp.MatchString(raws):
			authorization.Response = responseRegexp.ReplaceAllString(raws, "")

		case algorithmRegexp.MatchString(raws):
			authorization.Algorithm = algorithmRegexp.ReplaceAllString(raws, "")
		}
	}
	return authorization.Validator()
}
func (authorization *Authorization) Validator() error {
	if reflect.DeepEqual(nil, authorization) {
		return errors.New("authorization caller is not allowed to be nil")
	}
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if !authSchemaRegexp.MatchString(authorization.AuthSchema) {
		return errors.New("the value of the authschema field must be Digest")
	}
	if len(strings.TrimSpace(authorization.UserName)) == 0 {
		return errors.New("the username field is not allowed to be empty")
	}
	if len(strings.TrimSpace(authorization.Realm)) == 0 {
		return errors.New("the realm field is not allowed to be empty")
	}
	if len(strings.TrimSpace(authorization.Nonce)) == 0 {
		return errors.New("the nonce field is not allowed to be empty")
	}
	if err := authorization.Uri.Validator(); err != nil {
		return err
	}
	if len(strings.TrimSpace(authorization.Response)) == 0 {
		return errors.New("the response field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(md5)`).MatchString(authorization.Algorithm) {
		return errors.New("the value of the algorithm field must be MD5")
	}
	return nil
}
