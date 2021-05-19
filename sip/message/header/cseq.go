package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kokutas/gb28181/sip/lib"
)

type CSeq struct {
	sequenceNumber uint64 // sequence Number
	method         string // method
}

func (cseq *CSeq) SetSequenceNumber(number uint64) {
	cseq.sequenceNumber = number
}
func (cseq *CSeq) GetSequenceNumber() uint64 {
	return cseq.sequenceNumber
}

func (cseq *CSeq) SetMethod(method string) {
	cseq.method = method
}
func (cseq *CSeq) GetMethod() string {
	return cseq.method
}
func NewCSeq(number uint64, method string) *CSeq {
	return &CSeq{
		sequenceNumber: number,
		method:         method,
	}
}
func (cseq *CSeq) Raw() (string, error) {
	result := ""
	if err := cseq.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("CSeq: %d %s", cseq.sequenceNumber, strings.ToUpper(cseq.method))
	result += "\r\n"
	return result, nil
}
func (cseq *CSeq) Parse(raw string) error {
	if reflect.DeepEqual(nil, cseq) {
		return errors.New("cseq caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// cseq field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(cseq).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a cseq header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range lib.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")$"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if !methodsRegexp.MatchString(raw) {
		return errors.New("the value of the method field cannot be matched")
	}
	cseq.method = strings.ToUpper(strings.TrimSpace(methodsRegexp.FindString(raw)))
	raw = methodsRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// sequence Number regexp
	sequenceNumberRegexp := regexp.MustCompile(`\d+`)
	if sequenceNumberRegexp.MatchString(raw) {
		numberStr := sequenceNumberRegexp.FindString(raw)
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return err
		}
		cseq.sequenceNumber = uint64(number)
	}
	return nil
}

func (cseq *CSeq) Validator() error {
	if reflect.DeepEqual(nil, cseq) {
		return errors.New("cseq caller is not allowed to be nil")
	}
	if len(strings.TrimSpace(cseq.method)) == 0 {
		return errors.New("the method field is not allowed to be empty")
	}
	// methods regexp
	methodsRegexpStr := `(?i)(`
	for _, v := range lib.Methods {
		methodsRegexpStr += v + "|"
	}
	methodsRegexpStr = strings.TrimSuffix(methodsRegexpStr, "|")
	methodsRegexpStr += ")$"
	methodsRegexp := regexp.MustCompile(methodsRegexpStr)
	if !methodsRegexp.MatchString(cseq.method) {
		return errors.New("the value of the method field cannot be matched")
	}
	return nil
}

func (cseq *CSeq) String() string {
	result := ""
	result += fmt.Sprintf("%d", cseq.sequenceNumber)
	if len(strings.TrimSpace(cseq.method)) > 0 {
		result += fmt.Sprintf(" %s", cseq.method)
	}
	return result
}
