package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewCSeq(t *testing.T) {
	cseq := NewCSeq(1, "register")
	fmt.Println(cseq.GetMethod())
}

func TestCSeq_Raw(t *testing.T) {
	cseqs := []*CSeq{
		NewCSeq(1, "bye"),
		NewCSeq(1, "invite"),
	}
	for _, cseq := range cseqs {
		fmt.Print(cseq.Raw())
	}
}

func TestCSeq_Parse(t *testing.T) {
	raws := []string{
		"CSeq: 1 invite",
		"cseq: 0 register",
	}
	for _, raw := range raws {
		cseq := new(CSeq)
		if err := cseq.Parse(raw); err != nil {
			log.Fatal(err, raw)
		}
		fmt.Print(cseq.Raw())
	}
}
