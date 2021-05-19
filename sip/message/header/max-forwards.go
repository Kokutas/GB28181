package header

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type MaxForwards struct {
	forwards uint8 // forwards
}

func (maxForwards *MaxForwards) SetForwards(forwards uint8) {
	maxForwards.forwards = forwards
}
func (maxForwards *MaxForwards) GetForwards() uint8 {
	return maxForwards.forwards
}
func NewMaxForwards(forwards uint8) *MaxForwards {
	return &MaxForwards{
		forwards: forwards,
	}
}

func (maxForwards *MaxForwards) Raw() (string, error) {
	result := ""
	if err := maxForwards.Validator(); err != nil {
		return result, err
	}
	result += fmt.Sprintf("Max-Forwards: %d", maxForwards.forwards)
	result += "\r\n"
	return result, nil
}
func (maxForwards *MaxForwards) Parse(raw string) error {
	if reflect.DeepEqual(nil, maxForwards) {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	raw = regexp.MustCompile(`\r`).ReplaceAllString(raw, "")
	raw = regexp.MustCompile(`\n`).ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	if len(strings.TrimSpace(raw)) == 0 {
		return errors.New("the raw parameter is not allowed to be empty")
	}
	// max-forwards field regexp
	fieldRegexp := regexp.MustCompile(`(?i)(max-forwards).*?:`)
	if !fieldRegexp.MatchString(raw) {
		return errors.New("raw is not a max-forwards header field")
	}
	raw = fieldRegexp.ReplaceAllString(raw, "")
	raw = strings.TrimPrefix(raw, " ")
	raw = strings.TrimSuffix(raw, " ")
	// forwards regexp
	forwardsRegexp := regexp.MustCompile(`\d+`)
	if forwardsRegexp.MatchString(raw) {
		forwardStr := forwardsRegexp.FindString(raw)
		forwards, err := strconv.Atoi(forwardStr)
		if err != nil {
			return err
		}
		maxForwards.forwards = uint8(forwards)
	}
	return maxForwards.Validator()
}
func (maxForwards *MaxForwards) Validator() error {
	if reflect.DeepEqual(nil, maxForwards) {
		return errors.New("max-forwards caller is not allowed to be nil")
	}
	return nil
}
func (maxForwards *MaxForwards) String() string {
	result := ""
	result += fmt.Sprintf("%d", maxForwards.forwards)
	return result
}
