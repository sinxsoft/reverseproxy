package proxypass

import (
	"net/http"
	"strings"
)

var ConfigSeparator = "||"

func DoProxyPass(pp []string, r *http.Request) (bool, string, string) {
	for _, val := range pp {
		// example: /abc/def    和   /abc/def||  goURL 相比，是相同的
		if strings.HasPrefix(val, r.RequestURI+ConfigSeparator) {
			remote := strings.Split(val, ConfigSeparator)[1]
			return true, strings.TrimSpace(remote), strings.Split(val, ConfigSeparator)[2]
		}
	}
	return false, "", ""
}
