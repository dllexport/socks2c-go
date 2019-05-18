// +build darwin

package systemproxy

import "strings"

func getDefaultInterfaceName() string {
	str := execAndGetRes("networksetup", "-listallnetworkservices")
	res := strings.Split(str, "\n")
	return res[1]
}

func enableNoProxyimpl() {
	execAndGetRes("networksetup", "-setsocksfirewallproxystate", getDefaultInterfaceName(), "off")
	execAndGetRes("networksetup", "-setautoproxystate", getDefaultInterfaceName(), "off")
}

func enablePacImpl(url string) {
	execAndGetRes("networksetup", "-setautoproxyurl", getDefaultInterfaceName(), url)
}

func enableGlobalImpl(ip, port string) {

	execAndGetRes("networksetup", "-setsocksfirewallproxy", getDefaultInterfaceName(), ip, port)
}
