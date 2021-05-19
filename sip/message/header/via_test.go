package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewVia(t *testing.T) {
	vias := []*Via{
		NewVia("sip", 2.0, "udp", "192.168.0.1", 0, 0, "z9hG4bK", ""),
		NewVia("sip", 2.0, "udp", "3402000000", 0, 0, "z9hG4bK", ""),
	}
	for _, via := range vias {
		data, err := json.Marshal(via)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\r\n", data)
	}
}

func TestVia_Raw(t *testing.T) {
	vias := []*Via{
		NewVia("sip", 2.0, "udp", "192.168.0.1", 5060, 1, "z9hG4bK", ""),
		NewVia("sip", 2.0, "udp", "3402000000", 0, 0, "z9hG4bK", ""),
	}
	for _, via := range vias {
		str, err := via.Raw()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(str)
	}
}

func TestVia_Parse(t *testing.T) {
	raws := []string{
		"Via: SIP/2.0/UDP 192.168.0.1:5060;rport;branch=z9hG4bK",
		"Via: SIP/2.0/UDP 3402000000;rport=5060;branch=z9hG4bK;received=123.234.123.5",
	}
	for _, raw := range raws {
		via := new(Via)
		if err := via.Parse(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Print(via.Raw())

	}
}
