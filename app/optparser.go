package app

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"../systemproxy"
	"./config"
	"./logger"
)

func Parse() {

	optKey := flag.String("k", "", "key for the proxy connection")
	optLog := flag.Int("log", 0, "set the log level, the higher, the more details")

	optServerHost := flag.String("s", "", "server endpoint")
	optSocks5Host := flag.String("c", "127.0.0.1:1080", "local socks5 server endpoint")

	optVersion := flag.Bool("v", false, "Version Infomation")
	optStop := flag.Bool("stop", false, "Stop socks2c that is currently running")

	optPac := flag.Bool("pac", false, "Enable pac mode")
	optGlobalProxy := flag.Bool("gp", false, "enable global proxy mode")

	flag.Parse()

	if *optVersion {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if *optStop {
		SendStopSingal()
		os.Exit(0)
	}

	if optLog != nil {
		logger.SetLogLevel(intabs(*optLog))
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
