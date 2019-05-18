//go:generate statik -src=./assets
//go:generate go fmt systemproxy/statik.go
package systemproxy

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "../statik"
	"github.com/rakyll/statik/fs"
)

func EnablePac() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	_, err = statikFS.Open("/proxy.pac")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	http.Handle("/asset/", http.StripPrefix("/asset", http.FileServer(statikFS)))
	go http.ListenAndServe("127.0.0.1:65533", nil)

	execAndGetRes("networksetup", "-setautoproxyurl", getDefaultInterfaceName(), "http://127.0.0.1:65533/asset/proxy.pac")

}
