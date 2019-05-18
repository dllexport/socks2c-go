package systemproxy

func EnableGlobal(ip, port string) {

	execAndGetRes("networksetup", "-setsocksfirewallproxy", getDefaultInterfaceName(), ip, port)
}
