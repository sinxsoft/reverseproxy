package main

import (
	rr "github.com/sinxsoft/reverseproxy/roundrobin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var RR = rr.NewWeightedRR(rr.RR_NGINX)

type handle struct {
	addrs []string
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := RR.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func startServer() {
	//被代理的服务器host和port
	h := &handle{}
	h.addrs = []string{"127.0.0.1:9001", "127.0.0.1:9001"}

	w := 1
	for _, e := range h.addrs {
		RR.Add(e, w)
		w++
	}

	go func() {
		err := http.ListenAndServe(":8888", h)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}()

	go func() {
		err := http.ListenAndServeTLS(":443", "/Users/henrik/Documents/golang/src/golang-ReverseProxy/server/cacert.pem",
			"/Users/henrik/Documents/golang/src/golang-ReverseProxy/server/privkey.pem", h)
		if err != nil {
			log.Fatalln("ListenAndServeTLS: ", err)
		}
	}()

}

func main() {
	startServer()
	<-Channel
}

var Channel chan (bool)
