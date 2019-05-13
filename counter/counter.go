package counter

var TCP_PROXY_COUNT uint64 = 0
var UDP_PROXY_COUNT uint64 = 0

type ProxyCount struct {
	tcp uint64
	udp uint64
}

func Get() ProxyCount {
	return ProxyCount{TCP_PROXY_COUNT, UDP_PROXY_COUNT}
}
