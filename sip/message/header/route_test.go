package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewRoute(t *testing.T) {
	route := NewRoute("",
		NewUri("sip", "34020000001320000001", "192.168.0.1", 5060, map[string]interface{}{"lr": ""}),
		NewUri("sip", "34020000001320000001", "3402000000", 0, map[string]interface{}{"lr": ""}),
	)
	data, err := json.Marshal(route)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestRoute_Raw(t *testing.T) {
	route := NewRoute("",
		NewUri("sip", "34020000001320000001", "192.168.0.1", 5060, map[string]interface{}{"lr": ""}),
		NewUri("sip", "34020000001320000001", "3402000000", 0, map[string]interface{}{"lr": ""}),
	)
	raw, err := route.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(raw)
}

func TestRoute_Parse(t *testing.T) {
	raw := "Route: <sip:34020000001320000001@192.168.0.1:5060;lr>, <sip:34020000001320000001@3402000000;lr>\r\n"
	route := new(Route)
	if err := route.Parse(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(route.Raw())
}
