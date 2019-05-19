// +build linux

package systemproxy


func enableNoProxyimpl() {
	execAndGetRes("gsettings","set", "org.gnome.system.proxy", "mode", "'none'")
}

func enablePacImpl(url string) {
	execAndGetRes("gsettings","set", "org.gnome.system.proxy", "mode", "'auto'")
	execAndGetRes("gsettings", "set", "org.gnome.system.proxy", "autoconfig-url", url)
}

func enableGlobalImpl(ip, port string) {
	execAndGetRes("gsettings","set", "org.gnome.system.proxy", "host", ip)
	execAndGetRes("gsettings","set", "org.gnome.system.proxy", "port", port)
	execAndGetRes("gsettings","set", "org.gnome.system.proxy", "mode", "'manual'")
}
