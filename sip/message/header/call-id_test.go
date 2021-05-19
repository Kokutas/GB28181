package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewCallID(t *testing.T) {
	callId := NewCallID("140a92f15c94d76d62a4fcd2d3558000", "192.168.0.26")
	data, err := json.Marshal(callId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestCallID_Raw(t *testing.T) {
	callId := NewCallID("140a92f15c94d76d62a4fcd2d3558000", "")
	str, err := callId.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(str)
}

func TestCallID_Parse(t *testing.T) {
	raws := []string{
		"Call-ID: 140a92f15c94d76d62a4fcd2d3558000",
		"Call-IDs: 140a92f15c94d76d62a4fcd2d3558000@192.168.0.26:5060",
	}
	for _, raw := range raws {
		callId := new(CallID)
		if err := callId.Parse(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Print(callId.Raw())
	}
}
