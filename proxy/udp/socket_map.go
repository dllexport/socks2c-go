package udp

import "net"
import "sync"

var local_socket *net.UDPConn = nil

func SetLocal(s *net.UDPConn) {
	local_socket = s
}

func GetLocal() *net.UDPConn {
	return local_socket
}

type socket_map_safe struct {
	sync.RWMutex
	Map map[string]*net.UDPConn
}

func CreateSocketMap() *socket_map_safe {
	sm := new(socket_map_safe)
	sm.Map = make(map[string]*net.UDPConn)
	return sm
}

// if key not exist, return nil
func (sm *socket_map_safe) Read(key string) *net.UDPConn {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}
func (sm *socket_map_safe) Write(key string, value *net.UDPConn) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}

func (sm *socket_map_safe) Delete(key string) {
	sm.Lock()
	delete(sm.Map, key)
	sm.Unlock()
}

var socket_map = CreateSocketMap()
