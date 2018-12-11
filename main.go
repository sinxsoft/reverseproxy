package main

import (
	"fmt"
	"github.com/sinxsoft/reverseproxy/config"
	"github.com/sinxsoft/reverseproxy/filter"
	hd "github.com/sinxsoft/reverseproxy/handler"
	"log"
	"net/http"
	"strings"
)

func startServer(cfg config.Config) {
	//被代理的服务器host和port
	h := &hd.HttpHandler{}
	//h.addrs = []string{"127.0.0.1:9001", "127.0.0.1:9001"}
	h.Addrs = strings.Split(cfg.RoundRobin.URLs, ",")
	w := 1
	for _, e := range h.Addrs {
		hd.RR.Add(e, w)
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
	cfg := config.GetConfig()
	fmt.Print("配置情况如下：", cfg)
	config.ConfigInst = cfg

	filter.GetFilterList()
	startServer(cfg)
	if cfg.HttpsOpen || cfg.HttpOpen {
		<-waitChannel
	} else {
		fmt.Print("没有配置任何监听端口，应用退出！")
	}
}

var waitChannel chan bool
