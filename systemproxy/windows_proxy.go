// +build windows

package systemproxy

import (
	"strconv"

	"./socks5tohttp"

	"golang.org/x/sys/windows/registry"
)

func enableNoProxyimpl() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		println(err)
		return
	}
	defer k.Close()
	k.SetDWordValue("ProxyEnable", 0)
	k.DeleteValue("AutoConfigURL")
	k.DeleteValue("ProxyServer")
}

func enablePacImpl(url string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		println(err)
		return
	}
	defer k.Close()

	err = k.SetStringValue("AutoConfigURL", url)
	if err != nil {
		println(err)
		return
	}
}

func enableGlobalImpl(ip, port string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		println(err)
		return
	}
	defer k.Close()

	httpPort, _ := strconv.Atoi(port)

	httpPort++

	httpPortStr := strconv.Itoa(httpPort)

	err = k.SetStringValue("ProxyServer", ip+":"+httpPortStr)
	if err != nil {
		println(err)
		return
	}
	err = k.SetDWordValue("ProxyEnable", 1)
	if err != nil {
		println(err)
		return
	}

	go socks5tohttp.Start(httpPortStr)

}
