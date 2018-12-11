package rewrite

import (
	"github.com/sinxsoft/reverseproxy/config"
	"net/http"
	"regexp"
	"strings"
)

var ConfigSeparator = " "

func DoRewrite(rewrite []config.Rewrite, r *http.Request) (bool, string, string) {
	for _, val := range rewrite {
		//先确定rule是否和request匹配，匹配才进行下一步操作
		if strings.HasSuffix(val.RewriteRule, r.RequestURI+ConfigSeparator) {
			headName := getHeadName(val.RewriteCond)
			headValue := r.Header.Get(headName)
			condValueRight := strings.TrimSpace(strings.Split(val.RewriteCond, ConfigSeparator)[1])
			re := regexp.MustCompile(condValueRight)
			m := re.FindAllStringSubmatch(headValue, -1)
			//确定cond正则的值匹配
			if len(m) > 0 {
				ruleMiddle := strings.TrimSpace(strings.Split(val.RewriteRule, ConfigSeparator)[1])
				ruleRight := strings.TrimSpace(strings.Split(val.RewriteRule, ConfigSeparator)[2])
				remote := ruleMiddle
				return true, strings.TrimSpace(remote), ruleRight
			}
		}
	}
	return false, "", ""
}

func getHeadName(cond string) string {
	re := regexp.MustCompile(`(?s)\$\{(.*)\}`)
	m := re.FindAllStringSubmatch(cond, -1)
	return m[0][1]
}
