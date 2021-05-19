package socket

import (
	"fmt"
	"log"
	"net"
)

func NewUdpServer(ip net.IP, port uint16) {
	lAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port)))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for {

	}

}
