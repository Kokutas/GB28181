package sip

import (
	"net/http"

	"github.com/kokutas/gb28181/sip/lib"
	"github.com/kokutas/gb28181/sip/line"
)

type SipUas struct {
	ID      string           `json:"ID"`
	Realm   string           `json:"Realm"`
	Address []*SipUasAddress `json:"Address"`
}
type SipUasAddress struct {
	IP        string
	Port      uint16
	Transport string
}

func (uas *SipUas) Response(headerRaw string) (string, error) {
	liner := line.NewStatusLine("sip", 2.0, http.StatusOK, lib.Success[http.StatusOK])
	// 判断是否是register
	// register的from和to不变
	// contact 和route要变，判断是否有rport，有的话要给received和rport的值
	// TODO: 重构代码，修改错误类型，是客户端错误还是……
	// 断言
	// if sipError, ok := err.(*lib.SipError); ok {
	// 	liner.SetStatusCode(sipError.Code)
	// } else {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }
	return lineRaw + headerRaw, nil
}
