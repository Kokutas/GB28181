package line

import (
	"fmt"
	"log"
	"testing"
)

func TestNewRequestUri(t *testing.T) {
	reqUri := NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil)
	fmt.Println(reqUri.GetSchema(), reqUri.GetUser(), reqUri.GetHost(), reqUri.GetPort(), reqUri.GetExtension())
	reqUri.SetSchema("sips")
	reqUri.SetHost("192.168.0.1")
	reqUri.SetPort(5060)
	fmt.Println(reqUri.String())
}

func TestRequestUri_Raw(t *testing.T) {
	uris := []*RequestUri{
		NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil),
		NewRequestUri("sip", "34020000002000000001", "192.168.0.26", 5060, nil),
	}
	for _, reqUri := range uris {
		raw, err := reqUri.Raw()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(raw)
	}
}

func TestRequestUri_Parse(t *testing.T) {
	raws := []string{
		"sip:34020000002000000001@3402000000",
		"sip:34020000002000000001@192.168.0.26:5060",
	}
	for _, raw := range raws {
		reqUri := new(RequestUri)
		err := reqUri.Parse(raw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reqUri.String())
	}
}
