package handler

import (
	"github.com/sinxsoft/reverseproxy/config"
	"github.com/sinxsoft/reverseproxy/filter"
	"github.com/sinxsoft/reverseproxy/proxypass"
	"github.com/sinxsoft/reverseproxy/rewrite"
	rr "github.com/sinxsoft/reverseproxy/roundrobin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var RR = rr.NewWeightedRR(rr.RR_NGINX)
var filterList []filter.Filter
var ProxyPass []string
var Rewrite []config.Rewrite
var once sync.Once

func initConfig() {
	once.Do(func() {
		filterList = filter.GetFilterList()
		ProxyPass = config.ConfigInst.ProxyPass
		Rewrite = config.ConfigInst.Rewrite
	})
}

type HttpHandler struct {
	Addrs []string
}

func (this *HttpHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	initConfig()
	//filter，proxypass，rewrite 三个处理逐步开展

	isFilter := false
	for _, val := range filterList {
		if val.IsFilter(request.RequestURI) {
			val.Execute(w, request)
			isFilter = true
			return
		}
	}
	if isFilter {
		return
	}

	addr := RR.Next().(string)

	if result, remote, path := proxypass.DoProxyPass(ProxyPass, request); result {
		serveThisUrl(remote, path, w, request)
		return
	}

	if result, remote, path := rewrite.DoRewrite(Rewrite, request); result {
		serveThisUrl(remote, path, w, request)
		return
	}
	serveThisUrl(addr, "", w, request)
}

func serveThisUrl(addr, path string, w http.ResponseWriter, request *http.Request) {
	//addr := RR.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//改变request的path为向目标服务器发送
	if path != "" {
		request.URL.Path = path
	}
	proxy.ServeHTTP(w, request)
}
