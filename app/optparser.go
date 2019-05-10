package app

import (
	"fmt"
	"os"

	"github.com/pborman/getopt"
)

func Parse() (key, server_ep, socks5_ep string) {

	optKey := getopt.StringLong("k", 0, "", "proxy key")
	optServerHost := getopt.StringLong("s", 0, "", "server ep")
	optSocks5Host := getopt.StringLong("c", 0, "", "local socks5 ep")
	optHelp := getopt.BoolLong("help", 0, "Help")
	optVersion := getopt.BoolLong("v", 0, "Version Infomation")
	optStop := getopt.BoolLong("stop", 0, "stop socks2c-go")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if *optVersion {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}
	if *optStop {
		SendStopSingal()
		os.Exit(0)
	}

	return *optKey, *optServerHost, *optSocks5Host
}
