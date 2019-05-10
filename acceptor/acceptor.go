package acceptor

import (
	"fmt"
	"net"
	"os"
)
import tcp "../proxy"

var acceptor_ep string

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func Init(server, socks5 string) {

}

func Accept() {
	acceptor, err := net.Listen("tcp", acceptor_ep)
	checkError(err)

	for {
		conn, err := acceptor.Accept()

		if err != nil {
			continue
		}

		go tcp.HandleConnection(conn)
	}
}
