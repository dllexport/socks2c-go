package acceptor

import (
	"fmt"
	"net"
	"os"
)
import tcp "../proxy"
import "../protocol"

var acceptor_ep string
var remote_ep string

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func endpointCheck(endpoint string) {
	_, err := net.Dial("udp", endpoint)
	if err != nil {
		fmt.Printf("%s is not a vaild endpoint\n", endpoint)
		os.Exit(-1)
	}
}

func Init(server, socks5 string) {
	if len(server) == 0 {
		fmt.Printf("--s missing\n")
		os.Exit(-1)
	}
	if len(socks5) == 0 {
		fmt.Printf("--c missing\n")
		os.Exit(-1)
	}

	endpointCheck(socks5)
	endpointCheck(server)

	acceptor_ep = socks5
	remote_ep = server
}

func Accept() {
	acceptor, err := net.Listen("tcp", acceptor_ep)
	checkError(err)

	fmt.Printf("[Client] TcpProxy started, Server: [%s], Key: [%s], Local socks5 Port: [%s]\n", remote_ep, protocol.GetKey(), acceptor_ep)
	for {
		conn, err := acceptor.Accept()

		if err != nil {
			continue
		}

		go tcp.HandleConnection(conn, remote_ep)
	}
}
