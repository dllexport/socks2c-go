package app

import (
	"flag"
	"fmt"
	"os"

	"socks2c-go/app/config"
)

func Parse() {

	optKey := flag.String("k", "", "key for the proxy connection")

	optServerHost := flag.String("s", "", "server endpoint")
	optSocks5Host := flag.String("c", "127.0.0.1:1080", "local socks5 server endpoint")

	optVersion := flag.Bool("v", false, "Version Infomation")
	optStop := flag.Bool("stop", false, "Stop socks2c that is currently running")

	flag.Parse()

	if *optVersion {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if *optStop {
		SendStopSingal()
		os.Exit(0)
	}

	config.Init(*optKey, *optServerHost, *optSocks5Host)

	return
}

func intabs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
