package main

import (
	"bufio"
	"fmt"
	"os"

	"./acceptor"
	"./app"
	"./app/logger"
	"./counter"
	"./systemproxy"
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

	logger.LOG_INFO("[proxy statistic] tcp: %d udp:%d\n", counter.TCP_PROXY_COUNT, counter.UDP_PROXY_COUNT)

	systemproxy.EnableNoProxy()
}
