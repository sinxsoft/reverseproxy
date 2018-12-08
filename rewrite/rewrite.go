package rewrite

import (
	"github.com/sinxsoft/reverseproxy/config"
	"net/http"
	"strings"
)

func DoRewrite(rewrite []config.Rewrite, r *http.Request) bool {

	for _, val := range rewrite {
		//if strings.Replace(val.)
	}
	return true
}
func DoProxyPass(pp []string, r *http.Request) (bool, string) {
	for _, val := range pp {
		// example: /abc/def    和   /abc/edf  goURL 相比，是相同的
		if strings.HasSuffix(val, r.RequestURI+" ") {
			//serveThisUrl(addr,w,r)
			remote := strings.Split(val, " ")[1]
			return true, strings.TrimSpace(remote)
		}
	}
	return false, ""
}
