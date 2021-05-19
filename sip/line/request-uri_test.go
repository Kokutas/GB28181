package line

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewRequestUri(t *testing.T) {
	uri := NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil)
	data, err := json.Marshal(uri)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestRequestUri_Raw(t *testing.T) {
	uris := []*RequestUri{
		NewRequestUri("sip", "34020000002000000001", "3402000000", 0, nil),
		NewRequestUri("sip", "34020000002000000001", "192.168.0.26", 5060, nil),
	}
	for _, uri := range uris {
		raw, err := uri.Raw()
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
		uri := new(RequestUri)
		err := uri.Parse(raw)
		if err != nil {
			log.Fatal(err)
		}
		data, err := json.Marshal(uri)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\r\n", data)
	}
}
