package header

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNewMaxForwards(t *testing.T) {
	maxForwards := NewMaxForwards(70)
	data, err := json.Marshal(maxForwards)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\r\n", data)
}

func TestMaxForwards_Raw(t *testing.T) {
	maxForwards := NewMaxForwards(70)
	raw, err := maxForwards.Raw()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(raw)
}

func TestMaxForwards_Parse(t *testing.T) {
	raw := "Max-Forwards: 70"
	maxForwards := new(MaxForwards)
	if err := maxForwards.Parse(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(maxForwards.Raw())
}
