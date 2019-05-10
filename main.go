package main

import (
	"fmt"
	"os"
)
import "github.com/pborman/getopt"
import "./libsodium"
import "./acceptor"
import "./protocol"

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	optKey := getopt.StringLong("k", 0, "", "proxy key")
	optServerHost := getopt.StringLong("s", 0, "", "server ep")
	optSocks5Host := getopt.StringLong("c", 0, "", "local socks5 ep")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	protocol.SetKey(*optKey)

	acceptor.Init(*optServerHost, *optSocks5Host)

	libsodium.Init()

	acceptor.Accept()
}
