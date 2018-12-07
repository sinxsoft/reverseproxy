package config

import (
	"encoding/json"
	"io/ioutil"
)

var ConfigInst Config

type RoundRobin struct {
	URLs string //127.0.0.1:9001, 127.0.0.1:9001 逗号分隔
}

type Config struct {
	HttpPort   string
	HttpsPort  string
	CertFile   string
	KeyFile    string
	HttpOpen   bool
	HttpsOpen  bool
	RoundRobin RoundRobin
	Rewrite    []Rewrite //[]string `json:"rewriteCondition"`
	ProxyPass  []string  `json:"proxyPass"`
	//Filter     []string  `json:"filter"` //执行的脚本，lua,shell......
	UrlFilter []UrlFilter `json:"urlFilter"` //执行的脚本，lua,shell......
}

type Rewrite struct {
	RewriteCond string `json:"rewriteCond"`
	RewriteRule string `json:"rewriteRule"`
}

type UrlFilter struct {
	UrlPath     string `json:"urlPath"`
	ExecuteDesc string `json:"executeDesc"`
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
func GetConfig() Config {
	JsonParse := NewJsonStruct()
	v := Config{}
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load("/Users/henrik/Documents/golang/src/github.com/sinxsoft/reverseproxy/config/config.json", &v)
	return v
}
