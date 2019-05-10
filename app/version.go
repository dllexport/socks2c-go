package app

var v string = "0.0.1"
var version string = "socks2c-go " + v + " without UOUT"

func Version() string {
	return version
}
