package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
	//第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用标准库中的实例处理。
	//第二个参数，则是我们基于net/http标准库实现Web框架的入口
}

// handler echoes r.url.path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
}

// handler echoes r.url.header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
	}
}

//%q
//使用Go语法以及必须时使用转义，
//以双引号括起来的字符串或者字节切片[]byte，或者是以单引号括起来的数字
