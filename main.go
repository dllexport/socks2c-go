package main

import (
	"fmt"
	"os"
)

import "./acceptor"

import "./app"
import "./app/config"

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	key, server_ep, socks5_ep := app.Parse()

	app.SingleApp()

	config.Init(key, server_ep, socks5_ep)

	acceptor.Run()

}
