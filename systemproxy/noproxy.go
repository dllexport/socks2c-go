package systemproxy

func EnableNoProxy() {
	execAndGetRes("networksetup", "-setsocksfirewallproxystate", getDefaultInterfaceName(), "off")
	execAndGetRes("networksetup", "-setautoproxystate", getDefaultInterfaceName(), "off")
}
