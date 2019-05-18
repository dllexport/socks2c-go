package app

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"../systemproxy"
	"./config"
	"./logger"
	"github.com/pborman/getopt"
)

func Parse() {

	optKey := getopt.StringLong("k", 0, "", "proxy key")
	optLog := getopt.StringLong("log", 0, "0", "log level")
	optServerHost := getopt.StringLong("s", 0, "", "server ep")
	optSocks5Host := getopt.StringLong("c", 0, "", "local socks5 ep")
	optHelp := getopt.BoolLong("help", 0, "Help")
	optVersion := getopt.BoolLong("v", 0, "Version Infomation")
	optStop := getopt.BoolLong("stop", 0, "stop socks2c-go")
	optPac := getopt.BoolLong("pac", 0, "enable pac mode")
	optGlobalProxy := getopt.BoolLong("gp", 0, "enable global proxy mode")

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

	config.Init(*optKey, *optServerHost, *optSocks5Host)

	systemproxy.EnableNoProxy()

	if *optPac {
		logger.LOG_INFO("[Enable Pac]\n")
		systemproxy.EnablePac()
		return
	}

	if *optGlobalProxy {
		logger.LOG_INFO("[Enable Global Proxy]\n")
		_, err := net.Dial("udp", *optSocks5Host)
		if err != nil {
			fmt.Printf("%s is not a vaild endpoint\n", *optSocks5Host)
			os.Exit(-1)
		}
		res := strings.Split(*optSocks5Host, ":")
		systemproxy.EnableGlobal(res[0], res[1])
	}

	return
}

func intabs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
