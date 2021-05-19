package header

import (
	"fmt"
	"testing"
)

func TestNewExpires(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Println(expires.GetSeconds())
}

func TestExpires_Raw(t *testing.T) {
	expires := NewExpires(3600)
	fmt.Print(expires.Raw())
}

func TestExpires_Parse(t *testing.T) {
	raw := "expires: 5600"
	ex := new(Expires)
	fmt.Println(ex.Parse(raw))
	fmt.Print(ex.Raw())
}
