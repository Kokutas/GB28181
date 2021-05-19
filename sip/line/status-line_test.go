package line

import (
	"fmt"
	"log"
	"testing"
)

func TestNewStatusLine(t *testing.T) {
	statusLines := []*StatusLine{
		NewStatusLine("sip", 2.0, 200, "OK"),
		NewStatusLine("sip", 2.0, 200, "No"),
	}
	for _, statusLine := range statusLines {
		fmt.Print(statusLine.GetSchema(), statusLine.GetVersion(), statusLine.GetStatusCode(), statusLine.GetReasonPhrase())
	}
}

func TestStatusLine_Raw(t *testing.T) {
	statusLines := []*StatusLine{
		NewStatusLine("sip", 2.0, 200, "OK"),
		NewStatusLine("sip", 2.0, 200, "OK"),
	}
	for _, statusLine := range statusLines {
		str, err := statusLine.Raw()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(str)
	}
}

func TestStatusLine_Parse(t *testing.T) {
	raws := []string{
		"SIP/2.0 200 OK",
		"SIP/2.0 200 NOT found",
	}
	for _, raw := range raws {
		statusLine := new(StatusLine)
		if err := statusLine.Parse(raw); err != nil {
			log.Fatal(err)
		}
		fmt.Print(statusLine.Raw())
	}

}
