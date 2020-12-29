package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"socks2c-go/acceptor"
	"socks2c-go/app"
	"socks2c-go/counter"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	app.Parse()

	app.SingleApp()

	acceptor.Run()

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	log.Printf("[proxy statistic] tcp: %d udp:%d\n", counter.TCP_PROXY_COUNT, counter.UDP_PROXY_COUNT)

}
