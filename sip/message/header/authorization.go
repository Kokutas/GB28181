package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Authorization struct {
	authSchema string // auth-schema: Basic / Digest
	username   string // username
	realm      string // realm
	nonce      string // nonce
	uri        *Uri   // Uri
	response   string // response
	algorithm  string // algorithm
}

func (authorization *Authorization) SetAuthSchema(authSchema string) {
	authorization.authSchema = authSchema
}
func (authorization *Authorization) GetAuthSchema() string {
	return authorization.authSchema
}
func (authorization *Authorization) SetUserName(username string) {
	authorization.username = username
}
func (authorization *Authorization) GetUserName() string {
	return authorization.username
}
func (authorization *Authorization) SetRealm(realm string) {
	authorization.realm = realm
}
func (authorization *Authorization) GetRealm() string {
	return authorization.realm
}
func (authorization *Authorization) SetNonce(nonce string) {
	authorization.nonce = nonce
}
func (authorization *Authorization) GetNonce() string {
	return authorization.nonce
}

func (authorization *Authorization) SetUri(uri *Uri) {
	authorization.uri = uri
}
func (authorization *Authorization) GetUri() *Uri {
	if authorization.uri != nil {
		return authorization.uri
	}
	return nil
}
func (authorization *Authorization) SetResponse(response string) {
	authorization.response = response
}
func (authorization *Authorization) GetResponse() string {
	return authorization.response
}
func (authorization *Authorization) SetAlgorithm(algorithm string) {
	authorization.algorithm = algorithm
}
func (authorization *Authorization) GetAlgorithm() string {
	return authorization.algorithm
}
func NewAuthorization(authSchema string, username string, realm string, nonce string, uri *Uri, response string, algorithm string) *Authorization {
	return &Authorization{
		authSchema: authSchema,
		username:   username,
		realm:      realm,
		nonce:      nonce,
		uri:        uri,
		response:   response,
		algorithm:  algorithm,
	}
}

func (authorization *Authorization) Raw() (string, error) {
	result := ""
	if err := authorization.Validator(); err != nil {
		return result, err
	}
	uriStr, err := authorization.uri.Raw()
	if err != nil {
		return result, err
	}
	result += fmt.Sprintf("Authorization: %s username=\"%s\",realm=\"%s\",nonce=\"%s\",uri=\"%s\",response=\"%s\",algorithm=%s",
		strings.Title(authorization.authSchema), authorization.username, authorization.realm, authorization.nonce, uriStr, authorization.response, strings.ToUpper(authorization.algorithm))
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
		authorization.authSchema = authSchemaRegexp.FindString(raw)
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
			authorization.username = usernameRegexp.ReplaceAllString(raws, "")
		case realmRegexp.MatchString(raws):
			authorization.realm = realmRegexp.ReplaceAllString(raws, "")
		case nonceRegexp.MatchString(raws):
			authorization.nonce = nonceRegexp.ReplaceAllString(raws, "")
		case uriRegexp.MatchString(raws):
			authorization.uri = new(Uri)
			uriStr := uriRegexp.ReplaceAllString(raws, "")
			uriStr = regexp.MustCompile(`"`).ReplaceAllString(uriStr, "")
			uriStr = regexp.MustCompile(`<`).ReplaceAllString(uriStr, "")
			uriStr = regexp.MustCompile(`>`).ReplaceAllString(uriStr, "")
			if err := authorization.uri.Parse(uriStr); err != nil {
				return err
			}
		case responseRegexp.MatchString(raws):
			authorization.response = responseRegexp.ReplaceAllString(raws, "")

		case algorithmRegexp.MatchString(raws):
			authorization.algorithm = algorithmRegexp.ReplaceAllString(raws, "")
		}
	}
	return authorization.Validator()
}
func (authorization *Authorization) Validator() error {
	if reflect.DeepEqual(nil, authorization) {
		return errors.New("authorization caller is not allowed to be nil")
	}
	authSchemaRegexp := regexp.MustCompile(`(?i)(digest|basic)`)
	if !authSchemaRegexp.MatchString(authorization.authSchema) {
		return errors.New("the value of the authschema field must be Digest")
	}
	if len(strings.TrimSpace(authorization.username)) == 0 {
		return errors.New("the username field is not allowed to be empty")
	}
	if len(strings.TrimSpace(authorization.realm)) == 0 {
		return errors.New("the realm field is not allowed to be empty")
	}
	if len(strings.TrimSpace(authorization.nonce)) == 0 {
		return errors.New("the nonce field is not allowed to be empty")
	}
	if err := authorization.uri.Validator(); err != nil {
		return err
	}
	if len(strings.TrimSpace(authorization.response)) == 0 {
		return errors.New("the response field is not allowed to be empty")
	}
	if !regexp.MustCompile(`(?i)(md5)`).MatchString(authorization.algorithm) {
		return errors.New("the value of the algorithm field must be MD5")
	}
	return nil
}
