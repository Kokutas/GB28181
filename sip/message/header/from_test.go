package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewFrom(t *testing.T) {
	from := NewFrom("", NewUri("sip", "34020000001320000001", "192.168.0.1", 0, nil), "")
	fmt.Println(from.GetDisplayName(), from.GetAddress().GetUser())
}

func TestFrom_Raw(t *testing.T) {
	from := NewFrom("3402", NewUri("sip", "34020000001320000001", "340200", 0, nil), "")
	str, err := from.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", str)
}

func TestFrom_Parse(t *testing.T) {
	raw := "From: \"3402\" sip:34020000001320000001@340200"
	from := new(From)
	fmt.Println(from.Parse(raw))
	fmt.Println(from.Raw())
}
