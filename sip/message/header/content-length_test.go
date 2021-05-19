package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewContentLength(t *testing.T) {
	cl := NewContentLength(0)
	data, err := json.Marshal(cl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
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
