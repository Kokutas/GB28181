package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Head struct {
	*Via
	*From
	*To
	*CallID
	*CSeq
	*Contact
	*MaxForwards
	*Expires
	*ContentLength
}

func NewHead(via *Via, from *From, to *To, callId *CallID, cseq *CSeq, contact *Contact, maxForwards *MaxForwards, expires *Expires, contentLength *ContentLength) *Head {
	return &Head{
		Via:           via,
		From:          from,
		To:            to,
		CallID:        callId,
		CSeq:          cseq,
		Contact:       contact,
		MaxForwards:   maxForwards,
		Expires:       expires,
		ContentLength: contentLength,
	}
}

func (head *Head) Raw() (string, error) {
	result := ""
	if err := head.Validator(); err != nil {
		return result, err
	}
	via, err := head.Via.Raw()
	if err != nil {
		return result, err
	}
	from, err := head.From.Raw()
	if err != nil {
		return result, err
	}
	to, err := head.To.Raw()
	if err != nil {
		return result, err
	}
	callId, err := head.CallID.Raw()
	if err != nil {
		return result, err
	}
	cseq, err := head.CSeq.Raw()
	if err != nil {
		return result, err
	}
	contact, err := head.Contact.Raw()
	if err != nil {
		return result, err
	}
	maxForwards, err := head.MaxForwards.Raw()
	if err != nil {
		return result, err
	}
	expires, err := head.Expires.Raw()
	if err != nil {
		return result, err
	}
	contentLength, err := head.ContentLength.Raw()
	if err != nil {
		return result, err
	}

	result += fmt.Sprintf("%s%s%s%s%s%s%s%s%s", via, from, to, callId, cseq, contact, maxForwards, expires, contentLength)
	result += "\r\n"
	return result, nil
}

func (head *Head) Parse(raw string) error {
	if reflect.DeepEqual(nil, head) {
		return errors.New("head caller is not allowed to be nil")
	}
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// regexp
	viaRegexp := regexp.MustCompile(`^(?i)(via).*?:.*`)
	fromRegexp := regexp.MustCompile(`^(?i)(from).*?:.*`)
	toRegexp := regexp.MustCompile(`^(?i)(to).*?:.*`)
	callIdRegexp := regexp.MustCompile(`^(?i)(call-id).*?:.*`)
	cseqRegexp := regexp.MustCompile(`^(?i)(cseq).*?:.*`)
	contactRegexp := regexp.MustCompile(`(?i)(contact).*?:.*`)
	maxForwardsRegexp := regexp.MustCompile(`(?i)(max-forwards).*?:.*`)
	expiresRegexp := regexp.MustCompile(`(?i)(expires).*?:.*`)
	contentLengthRegexp := regexp.MustCompile(`(?i)(content-length).*?:.*`)

	rawSlice := strings.Split(raw, "\n")
	for _, raws := range rawSlice {
		switch {
		case viaRegexp.MatchString(raws):
			head.Via = new(Via)
			if err := head.Via.Parse(raws); err != nil {
				return err
			}
		case fromRegexp.MatchString(raws):
			head.From = new(From)
			if err := head.From.Parse(raws); err != nil {
				return err
			}
		case toRegexp.MatchString(raws):
			head.To = new(To)
			if err := head.To.Parse(raws); err != nil {
				return err
			}
		case callIdRegexp.MatchString(raws):
			head.CallID = new(CallID)
			if err := head.CallID.Parse(raws); err != nil {
				return err
			}
		case cseqRegexp.MatchString(raws):
			head.CSeq = new(CSeq)
			if err := head.CSeq.Parse(raws); err != nil {
				return err
			}
		case contactRegexp.MatchString(raws):
			head.Contact = new(Contact)
			if err := head.Contact.Parse(raws); err != nil {
				return err
			}
		case maxForwardsRegexp.MatchString(raws):
			head.MaxForwards = new(MaxForwards)
			if err := head.MaxForwards.Parse(raws); err != nil {
				return err
			}
		case expiresRegexp.MatchString(raws):
			head.Expires = new(Expires)
			if err := head.Expires.Parse(raws); err != nil {
				return err
			}
		case contentLengthRegexp.MatchString(raws):
			head.ContentLength = new(ContentLength)
			if err := head.ContentLength.Parse(raws); err != nil {
				return err
			}
		}
	}

	return nil
}
func (head *Head) Validator() error {
	if reflect.DeepEqual(nil, head) {
		return errors.New("head caller is not allowed to be nil")
	}
	return nil
}
