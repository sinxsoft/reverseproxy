package handler

import (
	rr "github.com/sinxsoft/reverseproxy/roundrobin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var RR = rr.NewWeightedRR(rr.RR_NGINX)

type HttpHandler struct {
	Addrs []string
}

func (this *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	addr := RR.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
