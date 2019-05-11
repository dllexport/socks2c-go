package protocol

import (
	"fmt"
	"os"
	"unsafe"

	"../libsodium"
)

var proxy_key = [32]byte{0}
var key_str string

func SetKey(key string) {
	if len(key) == 0 {
		fmt.Printf("--k missing\n")
		os.Exit(-1)
	}
	key_str = key
	str := []byte(key)
	copy(proxy_key[0:32], str)
}

func GetKey() string {
	return key_str
}

type Protocol struct {
	NONCE          [12]byte
	LEN_TAG        [16]byte
	PAYLOAD_TAG    [16]byte
	PAYLOAD_LENGTH uint32
	PADDING_LENGTH uint32
}

func ProtocolSize() uint64 {
	return 52
}

func encryptHeader(protocol_hdr *Protocol, send_buff *[]byte) {

	encrypted_data, _, tag_out := libsodium.EncryptData(proxy_key, protocol_hdr.NONCE, (*send_buff)[44:], 8)

	copy((*send_buff)[44:], encrypted_data)
	copy(protocol_hdr.LEN_TAG[:], tag_out[:])

}

func encryptPayload(protocol_hdr *Protocol, send_buff *[]byte) {

	encrypted_data, _, tag_out := libsodium.EncryptData(proxy_key, protocol_hdr.NONCE, (*send_buff)[ProtocolSize():], uint64(protocol_hdr.PAYLOAD_LENGTH))

	copy((*send_buff)[ProtocolSize():], encrypted_data)
	copy(protocol_hdr.PAYLOAD_TAG[:], tag_out[:])

}

func decryptHeader(protocol_hdr *Protocol, encrypted_data *[]byte) bool {

	decrypted_data, decrypt_res := libsodium.DecryptData(proxy_key, protocol_hdr.NONCE, (*encrypted_data)[44:], 8, protocol_hdr.LEN_TAG)

	copy((*encrypted_data)[44:], decrypted_data)

	return decrypt_res

}

func decryptPayload(protocol_hdr *Protocol, encrypted_data *[]byte) bool {

	decrypted_data, decrypt_res := libsodium.DecryptData(proxy_key, protocol_hdr.NONCE, (*encrypted_data), uint64(protocol_hdr.PAYLOAD_LENGTH), protocol_hdr.PAYLOAD_TAG)

	if decrypt_res == true {
		copy(*(encrypted_data), decrypted_data)
	}

	return decrypt_res
}

func OnSocks5RequestSent(payload []byte) []byte {

	send_buff := make([]byte, ProtocolSize()+uint64(len(payload)))

	protocol_hdr := (*Protocol)(unsafe.Pointer(&send_buff[0]))

	copy((*protocol_hdr).NONCE[:], libsodium.RandomBytes(12))

	// we don't add obf here
	protocol_hdr.PAYLOAD_LENGTH = uint32(len(payload))
	protocol_hdr.PADDING_LENGTH = 0

	copy(send_buff[ProtocolSize():], payload)

	encryptPayload(protocol_hdr, &send_buff)
	encryptHeader(protocol_hdr, &send_buff)

	return send_buff
}

func OnPayloadReadFromLocal(payload []byte) []byte {
	return OnSocks5RequestSent(payload)
}

func OnPayloadHeaderReadFromRemote(protocol_hdr *Protocol, payload []byte) uint64 {

	decrypt_res := decryptHeader(protocol_hdr, &payload)

	if decrypt_res {
		// we still need to add PADDING even thought we don't do obf at client side
		return uint64(protocol_hdr.PAYLOAD_LENGTH + protocol_hdr.PADDING_LENGTH)
	}

	return 0
}

func OnPayloadReadFromRemote(protocol_hdr *Protocol, payload []byte) bool {

	decrypt_res := decryptPayload(protocol_hdr, &payload)

	return decrypt_res
}

func OnUdpPayloadReadFromClient(payload []byte) []byte {
	return OnSocks5RequestSent(payload)
}

func OnUdpPayloadReadFromRemote(payload []byte) ([]byte, bool) {

	protocol_hdr := (*Protocol)(unsafe.Pointer(&payload[0]))

	decrypt_res := decryptHeader(protocol_hdr, &payload)

	if decrypt_res == false {
		return nil, false
	}

	var payload_slice []byte = payload[ProtocolSize():]

	decrypt_res = decryptPayload(protocol_hdr, &payload_slice)

	if decrypt_res == false {
		return nil, false
	}

	return payload_slice[:protocol_hdr.PAYLOAD_LENGTH], true
}
