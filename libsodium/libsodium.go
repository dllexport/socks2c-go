package libsodium

import (
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS: /usr/local/lib/libsodium.a
// #include "sodium/crypto_aead_aes256gcm.h"
// #include "sodium/core.h"
// #include "sodium/randombytes.h"
import "C"

func Init() {
	fmt.Printf("[Init libsodium]\n")
	C.sodium_init()
}

func RandomBytes(size uint64) []byte {
	var data = make([]byte, size)
	C.randombytes_buf(unsafe.Pointer(&data[0]), C.ulong(len(data)))
	return data
}

func EncryptData(key [32]byte, nonce [12]byte, original_data []byte, original_data_length uint64) (encrypted_data []byte, tag_length uint64, tag_out [16]byte) {

	encrypted_data = make([]byte, original_data_length)

	tag_length = 0

	//fmt.Printf("before enc using:\noriginal_data: %v\nlen: %v\ntag_out:%v\nnonce: %v\nkey: %v\n", original_data, original_data_length, tag_out, nonce, key)

	original_data_copy := make([]byte, original_data_length)
	copy(original_data_copy, original_data)

	C.crypto_aead_aes256gcm_encrypt_detached(
		(*C.uchar)(unsafe.Pointer(&encrypted_data[0])),
		(*C.uchar)(unsafe.Pointer(&tag_out[0])),
		(*C.ulonglong)(unsafe.Pointer(&tag_length)),
		(*C.uchar)(unsafe.Pointer(&original_data_copy[0])), C.ulonglong(original_data_length),
		nil, 0, nil,
		(*C.uchar)(unsafe.Pointer(&nonce[0])),
		(*C.uchar)(unsafe.Pointer(&key)))

	return
}

func DecryptData(key [32]byte, nonce [12]byte, encrypted_data []byte, encrypted_data_length uint64, tag_in [16]byte) (decrypted_data []byte, res bool) {

	decrypted_data = make([]byte, encrypted_data_length)

	encrypted_data_copy := make([]byte, encrypted_data_length)
	copy(encrypted_data_copy, encrypted_data)

	//fmt.Printf("\n\n")

	//fmt.Printf("before dec using:\nencrypted_data: %v\ntag_out:%v\nnonce: %v\nkey: %v\n", encrypted_data, tag_in, nonce, key)

	dectypt_res := C.crypto_aead_aes256gcm_decrypt_detached(
		(*C.uchar)(unsafe.Pointer(&decrypted_data[0])),
		nil,
		(*C.uchar)(unsafe.Pointer(&encrypted_data_copy[0])),
		C.ulonglong(encrypted_data_length),
		(*C.uchar)(unsafe.Pointer(&tag_in[0])),
		nil,
		0,
		(*C.uchar)(unsafe.Pointer(&nonce[0])),
		(*C.uchar)(unsafe.Pointer(&key)))

	//fmt.Printf("%v\n\n", dectypt_res)

	if dectypt_res != 0 {
		return nil, false
	} else {
		return decrypted_data, true
	}

}
