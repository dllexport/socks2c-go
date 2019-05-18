package systemproxy

import "strings"

func getDefaultInterfaceNameImpl() string {
	str := execAndGetRes("networksetup", "-listallnetworkservices")
	res := strings.Split(str, "\n")
	return res[1]
}
