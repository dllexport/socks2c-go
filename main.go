package main

import (
	"fmt"
	"os"
)

import "./libsodium"
import "./acceptor"
import "./protocol"
import "./app"

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	key, server_ep, socks5_ep := app.Parse()

	app.SingleApp()

	protocol.SetKey(key)

	acceptor.Init(server_ep, socks5_ep)

	libsodium.Init()

	acceptor.Accept()
}
