package udp

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"socks2c-go/counter"

	protocol "socks2c-go/protocol"
	socks5 "socks2c-go/protocol/socks5"

	"socks2c-go/app/config"
)

var UDP_SESSION_TIMEOUT_DNS = 3 * time.Second
var UDP_SESSION_TIMEOUT_NORMAL = 30 * time.Second

// we set write SetWriteDeadline here
// cause whether a socket is timeout determined by the client side
// there's situation that client might send multiple packet while the remote server reply nothing
func setConnTimeout(conn net.Conn, port uint16) {
	if port == 53 || port == 5353 {
		conn.SetWriteDeadline(time.Now().Add(UDP_SESSION_TIMEOUT_DNS))
	} else {
		conn.SetWriteDeadline(time.Now().Add(UDP_SESSION_TIMEOUT_NORMAL))
	}
}

func HandlePacket(local_ep net.Addr, data []byte) {

	client_req := (*socks5.UDP_RELAY_PACKET)(unsafe.Pointer(&data[0]))

	if client_req.RSV != 0x00 {
		log.Printf("HandlePacket: RSV != 0x00 drop\n")
		return
	}
	if client_req.ATYP != 0x01 {
		log.Printf("HandlePacket: ATYP != 0x01 udp proxy support ipv4 only\n")
		return
	}

	ip, port, pass := socks5.ParseIpPortFromSocks5UdpPacket(data)

	if pass == false {
		return
	}

	log.Printf("[udp proxy] %s --> %s:%d\n", local_ep.String(), ip, port)

	conn, res := socket_map.Read(local_ep.String())

	if res == false {

		atomic.AddUint64(&counter.UDP_PROXY_COUNT, 1)

		//fmt.Printf("net udp connection\n")
		udpaddr, err := net.ResolveUDPAddr("udp4", config.ServerEndpoint)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		remote_conn, err := net.DialUDP("udp", nil, udpaddr)

		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		setConnTimeout(remote_conn, port)

		socket_map.Write(local_ep.String(), *remote_conn)

		sendToRemote(data, *remote_conn)

		go readFromRemote(local_ep, *remote_conn)

	} else {
		sendToRemote(data, conn)
	}

}

func closeRemoteSocket(local_ep net.Addr) {
	//socket_map.Read(local_ep.String()).Close()
	socket_map.Delete(local_ep.String())
}

func readFromRemote(local_ep net.Addr, conn net.UDPConn) {

	defer closeRemoteSocket(local_ep)

	var remote_recv_buff [1500]byte

	for {

		bytes_read, err := conn.Read(remote_recv_buff[:])

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
func sendToRemote(data []byte, conn net.UDPConn) {

	send_buff := protocol.OnUdpPayloadReadFromClient(data)

	_, err := conn.Write(send_buff)

	if err != nil {
		return
	}

	//fmt.Printf("send %d bytes to remote\n", byte_send)

}
