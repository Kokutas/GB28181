package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Header struct {
	*Authorization
	*CallID
	*Contact
	*ContentLength
	*ContentType
	*CSeq
	*Expires
	*From
	*MaxForwards
	*Route
	*To
	*UserAgent
	*Via
	*WWWAuthenticate
}

func NewHeader(
	authorization *Authorization,
	callId *CallID,
	contact *Contact,
	contentLength *ContentLength,
	contentType *ContentType,
	cSeq *CSeq,
	expires *Expires,
	from *From,
	maxForwards *MaxForwards,
	route *Route,
	to *To,
	userAgent *UserAgent,
	via *Via,
	wwwAuthenticate *WWWAuthenticate,

) *Header {
	return &Header{
		Authorization:   authorization,
		CallID:          callId,
		Contact:         contact,
		ContentLength:   contentLength,
		ContentType:     contentType,
		CSeq:            cSeq,
		Expires:         expires,
		From:            from,
		MaxForwards:     maxForwards,
		Route:           route,
		To:              to,
		UserAgent:       userAgent,
		Via:             via,
		WWWAuthenticate: wwwAuthenticate,
	}
}

func (head *Header) Raw() (string, error) {
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

func (head *Header) Parse(raw string) error {
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

	return head.Validator()
}
func (head *Header) Validator() error {
	if reflect.DeepEqual(nil, head) {
		return errors.New("head caller is not allowed to be nil")
	}
	// via,from,to,callid,contact,length,expires
	if !reflect.DeepEqual(nil, head.Authorization) {
		if err := head.Authorization.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.CallID) {
		if err := head.CallID.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.Contact) {
		if err := head.Contact.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.ContentLength) {
		if err := head.ContentLength.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.ContentType) {
		if err := head.ContentType.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.CSeq) {
		if err := head.CSeq.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.Expires) {
		if err := head.Expires.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.From) {
		if err := head.From.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.MaxForwards) {
		if err := head.MaxForwards.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.Route) {
		if err := head.Route.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.To) {
		if err := head.Route.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.UserAgent) {
		if err := head.UserAgent.Validator(); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(nil, head.Via) {
		if err := head.Via.Validator(); err != nil {
			return err
		}

	}
	if !reflect.DeepEqual(nil, head.WWWAuthenticate) {
		if err := head.WWWAuthenticate.Validator(); err != nil {
			return err
		}
	}
	return nil
}
func (head *Header) String() string {
	result := ""
	return result
}
