package header

import (
	"fmt"
	"testing"
)

func TestNewAuthorization(t *testing.T) {
	authorization := NewAuthorization("digest", "34020000001320000001", "3402000000", "nonce123", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), "response123", "md5")
	fmt.Println(authorization.GetAlgorithm())
}

func TestAuthorization_Raw(t *testing.T) {
	authorization := NewAuthorization("digest", "34020000001320000001", "3402000000", "nonce123", NewUri("sip", "34020000001320000001", "192.168.0.108", 5060, nil), "response123", "md5")

	fmt.Print(authorization.Raw())
}
