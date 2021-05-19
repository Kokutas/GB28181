package line

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewRequestLine(t *testing.T) {
	reqLine := NewRequestLine("registers", NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil), "sip", 2.0)
	data, err := json.Marshal(reqLine)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestRequestLine_Raw(t *testing.T) {
	reqLines := []*RequestLine{
		NewRequestLine("register", NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil), "sip", 2.0),
		NewRequestLine("register", NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil), "sip", 2.0),
		NewRequestLine("register", NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil), "sip", 2.0),
	}
	for _, reqLine := range reqLines {
		str, err := reqLine.Raw()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(str)
	}
}

func TestRequestLine_Parse(t *testing.T) {
	raws := []string{
		"REGISTER sip:34020000002000000001@3402000000 SIP/2.0",
		"REGISTER sip:34020000002000000001@3402000000 SIP/2.0",
		"REGISTER sip:34020000002000000001@3402000000 SIP/2.0",
	}
	for _, raw := range raws {
		reqLine := new(RequestLine)
		if err := reqLine.Parse(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Print(reqLine.Raw())
	}
}
