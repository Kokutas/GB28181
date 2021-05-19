package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewContentLength(t *testing.T) {
	cl := NewContentLength(0)
	fmt.Println(cl.GetLength())
}

func TestContentLength_Raw(t *testing.T) {
	cl := NewContentLength(10)
	fmt.Print(cl.Raw())
}

func TestContentLength_Parse(t *testing.T) {
	raws := []string{
		"Content-Length: 10\r\n",
		"Content-Length: 1000000000000000000000000000000\r\n",
	}
	for _, raw := range raws {
		cl := new(ContentLength)
		if err := cl.Parse(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Print(cl.Raw())
	}
}
