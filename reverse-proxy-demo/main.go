package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	address  string
	username string
	password string
	port     int
)

func init() {
	flag.StringVar(&address, "address", "http://localhost:8090", "target server address")
	flag.StringVar(&username, "username", "admin", "http basic auth username")
	flag.StringVar(&password, "password", "secret", "http basic auth password")
	flag.IntVar(&port, "port", 9002, "reverse proxy server port")
}

func main() {
	flag.Parse()

	fmt.Println(fmt.Sprintf("address: %s, username: %s, password: %s, port: %d", address, username, password, port))

	remote, err := url.Parse(address)
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("===> ", remote.Host, r.URL)
			r.Host = remote.Host
			r.SetBasicAuth(username, password)
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	log.Println("==> start reverse proxy server")
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
