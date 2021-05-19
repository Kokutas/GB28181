package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewWWWAuthenticate(t *testing.T) {
	wwwAuthenticate := NewWWWAuthenticate("digest", "3402000000", "nonce123", "md5")
	data, err := json.Marshal(wwwAuthenticate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestWWWAuthenticate_Raw(t *testing.T) {
	wwwAuthenticate := NewWWWAuthenticate("digest", "3402000000", "nonce123", "md5")
	raw, err := wwwAuthenticate.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(raw)
}

func TestWWWAuthenticate_Parse(t *testing.T) {
	raw := "WWW-Authenticate: Digest realm=\"3402000000\",nonce=\"nonce123\"\r\n"
	wwwAuthenticate := new(WWWAuthenticate)
	if err := wwwAuthenticate.Parse(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(wwwAuthenticate.Raw())
}
