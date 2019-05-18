package app

var v string = "1.1.3"
var version string = "socks2c-go " + v

func Version() string {
	return version
}
