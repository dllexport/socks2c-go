package main

import (
	"bufio"
	"fmt"
	"os"

	"./acceptor"
	"./app"
	"./app/config"
	"./counter"
)

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

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	fmt.Printf("[proxy statistic] tcp: %d udp:%d\n", counter.TCP_PROXY_COUNT, counter.UDP_PROXY_COUNT)
}
