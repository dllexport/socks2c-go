//go:generate statik -src=./assets
//go:generate go fmt systemproxy/statik.go
package systemproxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"../app/config"
	_ "../statik"
	"github.com/rakyll/statik/fs"
)

var PacData string

func enablePac(url string) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	file, err := statikFS.Open("/proxy.pac")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	pacContent := string(content)

	socks5Port := strings.Split(config.Socks5Endpoint, ":")

	PacData = strings.Replace(pacContent, "$PORT", socks5Port[1], 1)

	http.HandleFunc("/", pacReply)
	go http.ListenAndServe("127.0.0.1:65533", nil)

	enablePacImpl(url)
}

func pacReply(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, PacData)
}
