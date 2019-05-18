package tcp

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
	"unicode/utf8"
	"unsafe"

	"../../protocol"
	socks5 "../../protocol/socks5"
)
import logger "../../app/logger"
import "../../counter"

var remote_ep string

// if the tcp connection is idle for 60 , close it
var socket_timeout = 60 * time.Second

func handleMethodSelection(conn net.Conn) bool {

	socks5_method_buff := make([]byte, 3)

	_, err := io.ReadFull(conn, socks5_method_buff)

	if err != nil {
		return false
	}
	client_req := (*socks5.METHOD_REQ)(unsafe.Pointer(&socks5_method_buff[0]))

	if client_req.VER != 0x05 {
		logger.LOG_DEBUG("METHOD_REQ.VER != 0x05\n")
		logger.LOG_DEBUG("%v\n", socks5_method_buff[:3])
		return false
	}

	buff := make([]byte, 2)
	client_reply := (*socks5.METHOD_REPLY)(unsafe.Pointer(&buff[0]))
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

	client_req := (*socks5.SOCKS_REQ)(unsafe.Pointer(&req_buff[0]))

	if client_req.VER != 0x05 {
		return false
	}

	// we send socks reply back
	_, send_err := (*local_conn).Write(socks5.DEFAULT_SOCKS_REPLY[:])

	if send_err != nil {
		return false
	}

	switch client_req.ATYP {
	case socks5.IPV4:
		{
			ip, port, pass := socks5.ParseIpPortFromSocks5Request(req_buff)

			if pass == false {
				return false
			}

			logger.LOG_INFO("[tcp proxy] %s --> %s:%d\n", (*local_conn).RemoteAddr().String(), ip, port)

			break
		}
	case socks5.IPV6:
		{
			logger.LOG_INFO("ipv6 not support yet")
			return false
		}
	case socks5.DOMAINNAME:
		{
			domain, port, pass := socks5.ParseDomainPortFromSocks5Request(req_buff)

			if pass == false {
				return false
			}

			logger.LOG_INFO("[tcp proxy] %s:%d\n", domain, port)

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

		remote_conn.SetReadDeadline(time.Now().Add(socket_timeout))

		var protocol_hdr_buff = make([]byte, protocol.ProtocolSize())
		_, err := io.ReadFull(remote_conn, protocol_hdr_buff)

		if err != nil {
			logger.LOG_DEBUG("Read header err --> %s\n", err.Error())
			break
		}

		//fmt.Printf("read %v bytes hdr from remote\n", byte_read)

		protocol_hdr := (*protocol.Protocol)(unsafe.Pointer(&protocol_hdr_buff[0]))

		payload_len := protocol.OnPayloadHeaderReadFromRemote(protocol_hdr, protocol_hdr_buff)

		if payload_len == 0 {
			fmt.Printf("OnPayloadHeaderReadFromRemote err")
			break
		}

		logger.LOG_DEBUG("payload len %d\n", payload_len)

		remote_recv_buff := make([]byte, payload_len)

		byte_read, err := io.ReadFull(remote_conn, remote_recv_buff)

		if err != nil {
			logger.LOG_DEBUG("io.ReadFull err --> %v\n", err.Error())
			return
		}

		logger.LOG_DEBUG("read %d bytes payload from remote\n", byte_read)

		read_err := protocol.OnPayloadReadFromRemote(protocol_hdr, remote_recv_buff)

		if read_err != true {
			logger.LOG_DEBUG("[tcp proxy] decrypt err\n")
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

	atomic.AddUint64(&counter.TCP_PROXY_COUNT, 1)

	handleTunnelFlow(local_conn, remote_conn)
}
