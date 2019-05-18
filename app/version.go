package app

var v string = "1.1.0"
var version string = "socks2c-go " + v

func Version() string {
	return version
}
