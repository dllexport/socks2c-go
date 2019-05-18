package systemproxy

var default_pac_url string = "http://127.0.0.1:65533/asset/proxy.pac"

func EnableGlobal(ip, port string) {
	enableGlobalImpl(ip, port)
}

func EnableNoProxy() {
	enableNoProxyimpl()
}

func EnablePac() {
	enablePac(default_pac_url)
}
