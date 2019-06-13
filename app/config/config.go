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
		fmt.Printf("-s missing\n")
		os.Exit(-1)
	}
	if len(socks5) == 0 {
		fmt.Printf("-c missing\n")
		os.Exit(-1)
	}
	if len(key) == 0 {
		fmt.Printf("-k missing\n")
		os.Exit(-1)
	}
	EndpointCheck(socks5)
	EndpointCheck(server)

	Socks5Endpoint = socks5
	ServerEndpoint = server
	//Key = key

	protocol.SetKey(key)

	libsodium.Init()

}

func EndpointCheck(endpoint string) {
	_, err := net.Dial("udp", endpoint)
	if err != nil {
		fmt.Printf("%s is not a vaild endpoint\n", endpoint)
		os.Exit(-1)
	}
}
