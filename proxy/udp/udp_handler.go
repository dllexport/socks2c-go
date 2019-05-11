package udp

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"unsafe"

	"../../protocol"

	"../../counter"

	socks5 "../../protocol/socks5"

	config "../../app/config"
)

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

	fmt.Printf("[udp proxy] %s:%d\n", ip, port)

	if socket_map.read(local_ep) == nil {

		atomic.AddUint64(&counter.UDP_PROXY_COUNT, 1)

		//fmt.Printf("net udp connection\n")
		udpaddr, err := net.ResolveUDPAddr("udp4", config.ServerEndpoint)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		remote_conn, err := net.DialUDP("udp", nil, udpaddr)

		socket_map.write(local_ep, remote_conn)

		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(-1)
		}

		go readFromRemote(local_ep)
	}

	sendToRemote(data, socket_map.read(local_ep))

}

func closeRemoteSocket(local_ep net.Addr) {
	socket_map.read(local_ep).Close()
	socket_map.write(local_ep, nil)
}

func readFromRemote(local_ep net.Addr) {

	defer closeRemoteSocket(local_ep)

	var remote_recv_buff [1500]byte

	for {
		bytes_read, err := socket_map.read(local_ep).Read(remote_recv_buff[:])

		if err != nil {
			return
		}

		//fmt.Printf("read %d bytes\n", bytes_read)

		send_buff, res := protocol.OnUdpPayloadReadFromRemote(remote_recv_buff[:bytes_read])

		if res == false {
			return
		}

		//fmt.Printf("read %d bytes\n", bytes_read)

		_, err = local_socket.WriteTo(send_buff, local_ep)
		if err != nil {
			return
		}
	}

}

func sendToRemote(data []byte, conn *net.UDPConn) {

	send_buff := protocol.OnUdpPayloadReadFromClient(data)

	_, err := conn.Write(send_buff)

	if err != nil {
		return
	}

	//fmt.Printf("send %d bytes to remote\n", byte_send)

}
