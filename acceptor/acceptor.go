package acceptor

import (
	"fmt"
	"net"
	"os"
)
import tcp "../proxy/tcp"
import udp "../proxy/udp"
import "../protocol"
import config "../app/config"

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func tcpAccept() {
	acceptor, err := net.Listen("tcp", config.Socks5Endpoint)
	checkError(err)

	fmt.Printf("[Client] TcpProxy started, Server: [%s], Key: [%s], Local socks5 Port: [%s]\n", config.ServerEndpoint, protocol.GetKey(), config.Socks5Endpoint)
	for {
		conn, err := acceptor.Accept()

		if err != nil {
			continue
		}

		go tcp.HandleConnection(conn, config.ServerEndpoint)
	}
}

func udpAccept() {
	udpaddr, _ := net.ResolveUDPAddr("udp4", config.Socks5Endpoint)
	acceptor, err := net.ListenUDP("udp", udpaddr)
	checkError(err)

	fmt.Printf("[Client] UdpProxy started, Server: [%s], Key: [%s], Local socks5 Port: [%s]\n", config.ServerEndpoint, protocol.GetKey(), config.Socks5Endpoint)

	udp.SetLocal(acceptor)

	var local_recv_buff [1500]byte
	for {
		bytes_read, local_ep, err := acceptor.ReadFrom(local_recv_buff[:])

		if err != nil {
			continue
		}

		go udp.HandlePacket(local_ep, local_recv_buff[:bytes_read])
	}
}

func Run() {
	go tcpAccept()
	udpAccept()
}
