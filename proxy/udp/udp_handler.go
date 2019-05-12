package udp

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"../../protocol"

	"../../counter"

	socks5 "../../protocol/socks5"

	config "../../app/config"
)

var UDP_SESSION_TIMEOUT_DNS = 3 * time.Second
var UDP_SESSION_TIMEOUT_NORMAL = 30 * time.Second

func setConnTimeout(conn net.Conn, port uint16) {
	if port == 53 || port == 5353 {
		conn.SetDeadline(time.Now().Add(UDP_SESSION_TIMEOUT_DNS))
	} else {
		conn.SetDeadline(time.Now().Add(UDP_SESSION_TIMEOUT_NORMAL))
	}
}

func HandlePacket(local_ep net.Addr, data []byte) {

	client_req := (*socks5.UDP_RELAY_PACKET)(unsafe.Pointer(&data[0]))

	if client_req.RSV != 0x00 {
		return
	}
	if client_req.ATYP != 0x01 {
		return
	}

	ip, port, pass := socks5.ParseIpPortFromSocks5UdpPacket(data)

	if pass == false {
		return
	}

	fmt.Printf("[udp proxy] from %s to %s:%d\n", local_ep.String(), ip, port)

	if socket_map.Read(local_ep.String()) == nil {

		atomic.AddUint64(&counter.UDP_PROXY_COUNT, 1)

		//fmt.Printf("net udp connection\n")
		udpaddr, err := net.ResolveUDPAddr("udp4", config.ServerEndpoint)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		remote_conn, err := net.DialUDP("udp", nil, udpaddr)

		setConnTimeout(remote_conn, port)

		socket_map.Write(local_ep.String(), remote_conn)

		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		go readFromRemote(local_ep)
	}
	sendToRemote(data, socket_map.Read(local_ep.String()))
}

func closeRemoteSocket(local_ep net.Addr) {
	socket_map.Read(local_ep.String()).Close()
	socket_map.Delete(local_ep.String())
}

func readFromRemote(local_ep net.Addr) {

	defer closeRemoteSocket(local_ep)

	var remote_recv_buff [1500]byte

	for {
		bytes_read, err := socket_map.Read(local_ep.String()).Read(remote_recv_buff[:])

		if err != nil {
			//fmt.Printf("remote socket err --> %s\n", err.Error())
			return
		}

		//fmt.Printf("read %d bytes\n", bytes_read)

		send_buff, res := protocol.OnUdpPayloadReadFromRemote(remote_recv_buff[:bytes_read])

		if res == false {
			return
		}

		//fmt.Printf("read %d bytes\n", bytes_read)

		_, err = GetLocal().WriteTo(send_buff, local_ep)
		if err != nil {
			return
		}
	}

}

// conn might be nil if readFromRemote for that specific ep returned
// because closeRemoteSocket() is defer called from readFromRemote
// which will clear the conn in the map
// as the operation is guarded by mutex so as long as conn != nil,
// it's safe to continue the sendToRemote even if it's removed afterward
func sendToRemote(data []byte, conn *net.UDPConn) {

	if conn == nil {
		return
	}

	send_buff := protocol.OnUdpPayloadReadFromClient(data)

	_, err := conn.Write(send_buff)

	if err != nil {
		return
	}

	//fmt.Printf("send %d bytes to remote\n", byte_send)

}
