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

func (this *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	initConfig()
	//filter，proxypass，rewrite 三个处理逐步开展

	isFilter := false
	for _, val := range filterList {
		if val.IsFilter(r.RequestURI) {
			val.Execute(w, r)
			isFilter = true
			return
		}
	}
	if isFilter {
		return
	}

	addr := RR.Next().(string)

	//for _, val := range ProxyPass {
	//	// example: /abc/def    和   /abc/edf  goURL 相比，是相同的
	//	if strings.HasSuffix(val, r.RequestURI+" ") {
	//		serveThisUrl(addr,w,r)
	//		break
	//	}
	//}
	if result, remote := proxypass.DoProxyPass(ProxyPass, r); result {
		serveThisUrl(addr+remote, w, r)
		return
	}

	if rewrite.DoRewrite(Rewrite, r) {
		return
	}

	//remote, err := url.Parse("http://" + addr)
	//if err != nil {
	//	panic(err)
	//}
	//proxy := httputil.NewSingleHostReverseProxy(remote)
	//proxy.ServeHTTP(w, r)
}

func serveThisUrl(addr string, w http.ResponseWriter, r *http.Request) {
	//addr := RR.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
