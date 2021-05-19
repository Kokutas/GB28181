package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewContact(t *testing.T) {
	contact := NewContact("", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), nil)
	data, err := json.Marshal(contact)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestContact_Raw(t *testing.T) {
	contacts := []*Contact{
		NewContact("", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), nil),
		NewContact("", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), map[string]interface{}{"expires": 3600}),
		NewContact("34020000001320000001", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), map[string]interface{}{"expires": 3600}),
	}
	for _, contact := range contacts {
		str, err := contact.Raw()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(str)
	}
}

func TestContact_Parse(t *testing.T) {
	raws := []string{
		"Contact: <sip:34020000001320000001@192.168.0.108:5060>",
		"Contact: <sip:34020000001320000001@192.168.0.108:5060;lr>;expires=3600",
		`Contact: "34020000001320000001" sip:34020000001320000001@192.168.0.108:5060;lr;expires=3600`,
		`Contact: "34020000001320000001" sip:34020000001320000001@13402000000;lr;expires=3600`,
	}
	for _, raw := range raws {
		contact := new(Contact)
		if err := contact.Parse(raw); err != nil {
			log.Fatal(err, raw)
		}
		fmt.Println(contact.Raw())
	}
}
