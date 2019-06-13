package app

var v string = "1.1.5"
var version string = "socks2c-go " + v

func Version() string {
	return version
}
