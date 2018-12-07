package filter

import (
	"github.com/sinxsoft/reverseproxy/config"
	_ "github.com/sinxsoft/reverseproxy/config"
	"net/http"
	"strings"
)

type Filter struct {
	path        string
	configDes   string
	executeDesc string
	executeFile []byte
}

func (this *Filter) execute(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(strings.ToLower(this.configDes), ".lua") {

	} else if strings.HasSuffix(strings.ToLower(this.configDes), ".sh") {

	}
}

func (this *Filter) IsFilter(url string) bool {
	if this.path == url {
		return true
	}
	return false
}

func GetFilterList() []Filter {
	len := len(config.ConfigInst.UrlFilter)
	filters := make([]Filter, len, len)
	for idx, value := range config.ConfigInst.UrlFilter {
		ft := Filter{}
		ft.configDes = value.ExecuteDesc
		ft.path = value.UrlPath
		(filters)[idx] = ft
	}
	return filters
}
