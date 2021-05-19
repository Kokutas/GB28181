package main

import (
	"fmt"

	"github.com/kokutas/gb28181/sip/message/header"
)

func main() {
	callId := header.NewCallID("12", "")
	fmt.Println(callId.GetId())
}
