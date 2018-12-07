package main

import (
	"fmt"
	"github.com/sinxsoft/reverseproxy/config"
	rr "github.com/sinxsoft/reverseproxy/roundrobin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var RR = rr.NewWeightedRR(rr.RR_NGINX)

type Handler struct {
	addrs []string
}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	addr := RR.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func startServer(cfg config.Config) {
	//被代理的服务器host和port
	h := &Handler{}
	//h.addrs = []string{"127.0.0.1:9001", "127.0.0.1:9001"}
	h.addrs = strings.Split(cfg.RoundRobin.URLs, ",")
	w := 1
	for _, e := range h.addrs {
		RR.Add(e, w)
		w++
	}

	if cfg.HttpOpen {
		go func() {
			err := http.ListenAndServe(":"+cfg.HttpPort, h)
			if err != nil {
				log.Fatalln("ListenAndServe: ", err)
			}
		}()
		fmt.Print(cfg.HttpPort + "端口开启。。。")
	}

	if cfg.HttpsOpen {
		go func() {
			err := http.ListenAndServeTLS(":"+cfg.HttpsPort, cfg.CertFile, cfg.KeyFile, h)
			if err != nil {
				log.Fatalln("ListenAndServeTLS: ", err)
			}
		}()
		fmt.Print(cfg.HttpsPort + "端口开启。。。")
	}

}

func main() {
	config := config.GetConfig()
	fmt.Print("配置情况如下：", config)
	startServer(config)

	if config.HttpsOpen || config.HttpOpen {
		<-waitChannel
	} else {
		fmt.Print("没有配置任何监听端口，应用退出！")
	}

}

var waitChannel chan (bool)
