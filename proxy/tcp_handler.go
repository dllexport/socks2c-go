package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"

	"../protocol"
)

var remote_ep string

var socket_timeout = 60 * time.Second

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

var DEFAULT_SOCKS_REPLY = [...]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10}

func to_string(b byte) string {
	return strconv.Itoa(int(b))
}

func parseIpPortFromSocks5Request(socks5_req_buff []byte) (ip string, port uint16, pass bool) {

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

func parseDomainPortFromSocks5Request(socks5_req_buff []byte) (domain string, port uint16, pass bool) {

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

func handleMethodSelection(conn net.Conn) bool {

	socks5_method_buff := make([]byte, 3)

	_, err := io.ReadFull(conn, socks5_method_buff)

	if err != nil {
		return false
	}
	client_req := (*METHOD_REQ)(unsafe.Pointer(&socks5_method_buff[0]))

	if client_req.VER != 0x05 {
		return false
	}

	buff := make([]byte, 2)
	client_reply := (*METHOD_REPLY)(unsafe.Pointer(&buff[0]))
	client_reply.VER = 0x05
	client_reply.METHOD = 0x00

	_, err = conn.Write(buff)
	if err != nil {
		return false
	}

	return true
}

func handleSocks5Request(local_conn, remote_conn *net.Conn) bool {
	req_buff := make([]byte, 256)
	byte_read, err := (*local_conn).Read(req_buff)

	if err != nil {
		return false
	}

	//fmt.Printf("read %d bytes socks5 request", bytes_read)

	client_req := (*SOCKS_REQ)(unsafe.Pointer(&req_buff[0]))

	if client_req.VER != 0x05 {
		return false
	}

	// we send socks reply back
	_, send_err := (*local_conn).Write(DEFAULT_SOCKS_REPLY[:])

	if send_err != nil {
		return false
	}

	switch client_req.ATYP {
	case IPV4:
		{
			ip, port, pass := parseIpPortFromSocks5Request(req_buff)

			if pass == false {
				return false
			}

			fmt.Printf("[tcp proxy] %s:%d\n", ip, port)
			break
		}
	case IPV6:
		{
			fmt.Printf("ipv6 not support yet")
			return false
		}
	case DOMAINNAME:
		{
			domain, port, pass := parseDomainPortFromSocks5Request(req_buff)

			if pass == false {
				return false
			}

			fmt.Printf("[tcp proxy] %s:%d\n", domain, port)

			addr := net.ParseIP(domain)
			if addr != nil {
				fmt.Printf("given ip address but mark domain request\n")
				break
			}

			err := checkDomain(domain)

			if err != nil {
				fmt.Printf("%v is not a valid domain\n", domain)
				return false
			}

			// addr, err := net.LookupIP(domain)

			// if err != nil {
			// 	return false
			// }

			// fmt.Printf("resolved %s\n", addr[0])
			break
		}
	}
	return connectAndSend(remote_conn, req_buff[:byte_read])
}

// checkDomain returns an error if the domain name is not valid
// See https://tools.ietf.org/html/rfc1034#section-3.5 and
// https://tools.ietf.org/html/rfc1123#section-2.
func checkDomain(name string) error {
	switch {
	case len(name) == 0:
		return nil // an empty domain name will result in a cookie without a domain restriction
	case len(name) > 255:
		return fmt.Errorf("cookie domain: name length is %d, can't exceed 255", len(name))
	}
	var l int
	for i := 0; i < len(name); i++ {
		b := name[i]
		if b == '.' {
			// check domain labels validity
			switch {
			case i == l:
				return fmt.Errorf("cookie domain: invalid character '%c' at offset %d: label can't begin with a period", b, i)
			case i-l > 63:
				return fmt.Errorf("cookie domain: byte length of label '%s' is %d, can't exceed 63", name[l:i], i-l)
			case name[l] == '-':
				return fmt.Errorf("cookie domain: label '%s' at offset %d begins with a hyphen", name[l:i], l)
			case name[i-1] == '-':
				return fmt.Errorf("cookie domain: label '%s' at offset %d ends with a hyphen", name[l:i], l)
			}
			l = i + 1
			continue
		}
		// test label character validity, note: tests are ordered by decreasing validity frequency
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			// show the printable unicode character starting at byte offset i
			c, _ := utf8.DecodeRuneInString(name[i:])
			if c == utf8.RuneError {
				return fmt.Errorf("cookie domain: invalid rune at offset %d", i)
			}
			return fmt.Errorf("cookie domain: invalid character '%c' at offset %d", c, i)
		}
	}
	// check top level domain validity
	switch {
	case l == len(name):
		return fmt.Errorf("cookie domain: missing top level domain, domain can't end with a period")
	case len(name)-l > 63:
		return fmt.Errorf("cookie domain: byte length of top level domain '%s' is %d, can't exceed 63", name[l:], len(name)-l)
	case name[l] == '-':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d begins with a hyphen", name[l:], l)
	case name[len(name)-1] == '-':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d ends with a hyphen", name[l:], l)
	case name[l] >= '0' && name[l] <= '9':
		return fmt.Errorf("cookie domain: top level domain '%s' at offset %d begins with a digit", name[l:], l)
	}
	return nil
}

// data should be payload only
func connectAndSend(remote_conn *net.Conn, data []byte) (res bool) {

	var err_conn error
	*remote_conn, err_conn = net.Dial("tcp", remote_ep)

	if err_conn != nil {
		return false
	}

	send_buff := protocol.OnSocks5RequestSent(data)

	_, err := (*remote_conn).Write(send_buff)

	if err != nil {
		return false
	}
	//fmt.Printf("send %d bytes to remote\n", byte_send)

	return true
}

func upStream(local_conn, remote_conn net.Conn) {

	defer local_conn.Close()
	defer remote_conn.Close()

	local_recv_buff := make([]byte, 1500-52)

	for {
		local_conn.SetDeadline(time.Now().Add(socket_timeout))

		byte_read, err := local_conn.Read(local_recv_buff)

		if err != nil {
			//fmt.Print(err.Error())
			return
		}

		//fmt.Printf("read %v bytes from local\n", byte_read)
		//fmt.Printf("%v\n", local_recv_buff[:byte_read])

		send_buff := protocol.OnPayloadReadFromLocal(local_recv_buff[:byte_read])

		_, err = remote_conn.Write(send_buff)

		if err != nil {
			return
		}
		//fmt.Printf("send %d bytes to remote\n", byte_send)

	}
}
func downStream(local_conn, remote_conn net.Conn) {

	defer local_conn.Close()
	defer remote_conn.Close()

	for {
		remote_conn.SetDeadline(time.Now().Add(socket_timeout))

		var protocol_hdr_buff = make([]byte, protocol.ProtocolSize())
		_, err := io.ReadFull(remote_conn, protocol_hdr_buff)

		if err != nil {
			//fmt.Print(err.Error())
			return
		}

		//fmt.Printf("read %v bytes hdr from remote\n", byte_read)

		protocol_hdr := (*protocol.Protocol)(unsafe.Pointer(&protocol_hdr_buff[0]))

		payload_len := protocol.OnPayloadHeaderReadFromRemote(protocol_hdr, protocol_hdr_buff)

		if payload_len == 0 {
			fmt.Printf("OnPayloadHeaderReadFromRemote err")
			break
		}

		remote_recv_buff := make([]byte, payload_len)

		_, err = io.ReadFull(remote_conn, remote_recv_buff)

		//fmt.Printf("read %v bytes payload from remote\n", byte_read)

		read_err := protocol.OnPayloadReadFromRemote(protocol_hdr, remote_recv_buff)

		if read_err != true {
			fmt.Printf("decrypt err\n")
			return
		}

		_, err = local_conn.Write(remote_recv_buff[:protocol_hdr.PAYLOAD_LENGTH])
		if err != nil {
			//fmt.Printf("local_conn write err\n")
			return
		}
	}
}
func handleTunnelFlow(local_conn, remote_conn net.Conn) {
	go upStream(local_conn, remote_conn)
	go downStream(local_conn, remote_conn)
}

func HandleConnection(conn net.Conn, remote string) {

	remote_ep = remote

	var local_conn = conn
	var remote_conn net.Conn

	if handleMethodSelection(local_conn) == false {
		return
	}
	if handleSocks5Request(&local_conn, &remote_conn) == false {
		return
	}

	handleTunnelFlow(local_conn, remote_conn)
}
