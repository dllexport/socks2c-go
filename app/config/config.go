package config

import (
	"fmt"
	"net"
	"os"
)
import "../../protocol"
import "../../libsodium"

var ServerEndpoint string
var Socks5Endpoint string
var Key string

func Init(key, server, socks5 string) {
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

	Socks5Endpoint = socks5
	ServerEndpoint = server
	//Key = key

	protocol.SetKey(key)

	libsodium.Init()

}

func endpointCheck(endpoint string) {
	_, err := net.Dial("udp", endpoint)
	if err != nil {
		fmt.Printf("%s is not a vaild endpoint\n", endpoint)
		os.Exit(-1)
	}
}
