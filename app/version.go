package app

var v string = "1.0.0"
var version string = "socks2c-go " + v + " without UOUT"

func Version() string {
	return version
}
