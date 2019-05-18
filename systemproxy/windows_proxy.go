// +build windows

package systemproxy

import (
	"strconv"

	"./powershell"
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
	k.DeleteValue("AutoConfigURL")
	k.DeleteValue("ProxyServer")
	k.SetDWordValue("ProxyEnable", 0)

	refreshedSetting()
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

	k2, err := registry.OpenKey(registry.CURRENT_USER, `Software\Policies\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		println(err.Error())
		return
	}
	defer k2.Close()

	err = k2.SetDWordValue("EnableAutoproxyResultCache", 0)
	if err != nil {
		println(err.Error())
		return
	}

	refreshedSetting()
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

	err = k.SetStringValue("ProxyOverride", "<local>")
	if err != nil {
		println(err)
		return
	}

	err = k.SetDWordValue("ProxyEnable", 1)
	if err != nil {
		println(err)
		return
	}

	refreshedSetting()

	go socks5tohttp.Start(httpPortStr)

}

func refreshedSetting() {
	ps := powershell.New()

	if ps == nil {
		return
	}

	_, _, err := ps.Execute(
		"function Reload-InternetOptions\n" +
			"{\n" +
			"$signature = @'\n" +
			"[DllImport(\"wininet.dll\", SetLastError = true, CharSet=CharSet.Auto)]\n" +
			"public static extern bool InternetSetOption(IntPtr hInternet, int dwOption, IntPtr lpBuffer, int dwBufferLength);\n" +
			"'@\n" +
			"$interopHelper = Add-Type -MemberDefinition $signature -Name MyInteropHelper -PassThru\n" +
			"$INTERNET_OPTION_SETTINGS_CHANGED = 39\n" +
			"$INTERNET_OPTION_REFRESH = 37\n" +
			"$result1 = $interopHelper::InternetSetOption(0, $INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)\n" +
			"$result2 = $interopHelper::InternetSetOption(0, $INTERNET_OPTION_REFRESH, 0, 0)\n" +
			"$result1 -and $result2\n" +
			"}\n" +
			"Reload-InternetOptions\n",
	)
	ps.Execute("")
	// fmt.Println(stdout)
	// fmt.Println(stderr)
	if err != nil {
		//fmt.Println(err)
	}
	//fmt.Println("system proxy refreshed")
}
