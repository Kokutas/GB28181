package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewTo(t *testing.T) {
	to := NewTo("", NewUri("sip", "34020000001320000001", "192.168.0.1", 0, nil), "")
	fmt.Println(to.GetAddress())
}

func TestTo_Raw(t *testing.T) {
	to := NewTo("3402", NewUri("sip", "34020000001320000001", "340200", 0, nil), "")
	str, err := to.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", str)
}

func TestTo_Parse(t *testing.T) {
	raws := []string{
		// "From: \"3402\" sip:34020000001320000001@340200",
		"To: \"3402\" sip:34020000001320000001@340200",
		"To: <sip:34020000001320000001@340200>;tag=123",
		"To: <sip:34020000001320000001@192.168.0.1>;tag=123",
	}
	for _, raw := range raws {
		to := new(To)
		fmt.Println(to.Parse(raw))
		fmt.Println(to.Raw())
		fmt.Println(to.GetTag())
	}

}
