package header

import (
	"fmt"
	"log"
	"testing"
)

func TestNewUserAgent(t *testing.T) {
	userAgent := NewUserAgent("SIP UAS V3.0.0.833566")
	fmt.Println(userAgent.GetServer())
}

func TestUserAgent_Raw(t *testing.T) {
	userAgent := NewUserAgent("SIP UAS V3.0.0.833566")
	fmt.Println(userAgent.Raw())
}

func TestUserAgent_Parse(t *testing.T) {
	raw := "User-Agent: SIP UAS V3.0.0.833566\r\n"
	userAgent := new(UserAgent)
	if err := userAgent.Parse(raw); err != nil {
		log.Fatal(err)
	}
	fmt.Print(userAgent.Raw())
}
