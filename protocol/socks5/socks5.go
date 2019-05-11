package protocol

import (
	"encoding/binary"
	"strconv"
)

type METHOD_REQ struct {
	VER, METHOD, METHODS byte
}

type METHOD_REPLY struct {
	VER, METHOD byte
}

const (
	IPV4       = 0x01
	DOMAINNAME = 0x03
	IPV6       = 0x04
)

type SOCKS_REQ struct {
	VER, CMD, RSV, ATYP byte
}

type UDP_RELAY_PACKET struct {
	RSV        int16
	FRAG, ATYP byte
}

var DEFAULT_SOCKS_REPLY = [...]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10}

func to_string(b byte) string {
	return strconv.Itoa(int(b))
}

func ParseIpPortFromSocks5Request(socks5_req_buff []byte) (ip string, port uint16, pass bool) {

	for i := 4; i < 7; i++ {
		ip += to_string(socks5_req_buff[i])
		ip += "."
	}
	ip += to_string(socks5_req_buff[7])

	port_buff := make([]byte, 2)
	port_buff[1] = socks5_req_buff[9]
	port_buff[0] = socks5_req_buff[8]
	port = binary.BigEndian.Uint16(port_buff)

	return ip, port, true
}

func ParseDomainPortFromSocks5Request(socks5_req_buff []byte) (domain string, port uint16, pass bool) {

	domain_length := int(socks5_req_buff[4])

	for i := 0; i < domain_length; i++ {
		domain += string(socks5_req_buff[i+5])
	}

	port_buff := make([]byte, 2)
	port_buff[1] = socks5_req_buff[domain_length+6]
	port_buff[0] = socks5_req_buff[domain_length+5]
	port = binary.BigEndian.Uint16(port_buff)

	return domain, port, true
}

func ParseIpPortFromSocks5UdpPacket(socks5_req_buff []byte) (ip string, port uint16, pass bool) {
	return ParseIpPortFromSocks5Request(socks5_req_buff)
}
