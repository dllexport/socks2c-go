// +build windows

package socks5tohttp

import (
	"io"
	"net"
	"net/http"
	"time"

	"../../app/config"
	"golang.org/x/net/proxy"
)

func handleHTTP(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	tp := http.Transport{
		Dial: dialer.Dial,
	}
	resp, err := tp.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func handleTunnel(w http.ResponseWriter, req *http.Request, dialer proxy.Dialer) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	srcConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	dstConn, err := dialer.Dial("tcp", req.Host)
	if err != nil {
		srcConn.Close()
		return
	}

	srcConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go transfer(dstConn, srcConn)
	go transfer(srcConn, dstConn)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()

	io.Copy(dst, src)
}

func serveHTTP(w http.ResponseWriter, req *http.Request) {
	d := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	dialer, _ := proxy.SOCKS5("tcp", config.Socks5Endpoint, nil, d)

	if req.Method == "CONNECT" {
		handleTunnel(w, req, dialer)
	} else {
		handleHTTP(w, req, dialer)
	}
}

func Start(httpPortStr string) {
	http.ListenAndServe("127.0.0.1:"+httpPortStr, http.HandlerFunc(serveHTTP))
}
