package app

import (
	"fmt"
	"net"
	"os"
)

var uniqueSocks2cEndpoint = "127.0.0.1:44444"

func handleRead(conn *net.UDPConn) {
	var buff = [10]byte{0}
	conn.Read(buff[:])

	fmt.Printf("recv signal quitting\n")

	os.Exit(-1)
}

func SingleApp() {
	udpaddr, _ := net.ResolveUDPAddr("udp4", uniqueSocks2cEndpoint)

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		fmt.Printf("another socks2c is running\n")
		os.Exit(0)
	}
	go handleRead(conn)
}

func SendStopSingal() {
	udpaddr, _ := net.ResolveUDPAddr("udp4", uniqueSocks2cEndpoint)
	udpconn, _ := net.DialUDP("udp", nil, udpaddr)
	udpconn.Write([]byte("stop\r\n"))
	fmt.Printf("stopping socks2c-go\n")
}
