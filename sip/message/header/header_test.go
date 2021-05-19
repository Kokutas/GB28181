package header

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"
)

func TestNewHead(t *testing.T) {
	// 	schema := "sip"
	// 	version := 2.0
	// 	transport := "udp"
	// 	sentByHost := "192.168.0.108"
	// 	sentByPort := 5060
	// 	rport := 0
	// 	branch := "z9hG4bKe8cd231f5ba018400e711fbfd459171b"
	// 	received := ""
	// 	user := "34020000001320000001"
	// 	displayName := "34020000001320000001"
	// 	fromTag := "fromTag"
	// 	toTag := "toTag"
	// 	callId := "123"
	// 	NewHead(
	// 		NewVia(schema, version, transport, sentByHost, uint16(sentByPort), rport, branch, received),
	// 		NewFrom(displayName, NewUri(schema, user, sentByHost, uint16(sentByPort), nil), fromTag),
	// 		NewTo(displayName, NewUri(schema, user, sentByHost, uint16(sentByPort), nil), toTag),
	// 		NewCallID(callId, sentByHost),
	// 	)

	// 	_, _ = http.Get("www.baidu.com")

	for i := 0; i < 6; i++ {
		rsp, _ := http.Get("http://www.baidu.com")
		_, _ = ioutil.ReadAll(rsp.Body)
	}
	fmt.Println(runtime.NumGoroutine())
}
