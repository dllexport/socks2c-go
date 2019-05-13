package app

import (
	"fmt"
	"os"
	"strconv"

	"./logger"
	"github.com/pborman/getopt"
)

func Parse() (key, server_ep, socks5_ep string) {

	optKey := getopt.StringLong("k", 0, "", "proxy key")
	optLog := getopt.StringLong("log", 0, "0", "log level")
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

	if optLog != nil {
		s, err := strconv.Atoi(*optLog)
		if err != nil {
			fmt.Printf("--log err, enable default level 0\n")
			logger.SetLogLevel(0)
		}
		logger.SetLogLevel(intabs(s))
	}

	return *optKey, *optServerHost, *optSocks5Host
}

func intabs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
