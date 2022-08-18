package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

// http.handler，不带ps httprouter.Params入参
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Blog:%s \nWechat:%s", "www.flysnow.org", "flysnow_org")
}

// httprouter.handler，ps httprouter.Params
//:作为通配符，:开头的为可变参数，通过httprouter.Params.ByName获取
func UserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//一个新类型，用于存储域名对应的路由
type HostSwitch map[string]http.Handler

//实现http.Handler接口，进行不同域名的路由分发
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//根据域名获取对应的Handler路由，然后调用处理（分发机制）
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/user/:name", UserInfo)
	//使用gorilla包提供的http.handler，需要转换成用router.Handler注册http.handler
	//ServeMux路由分发->调用中间件1->调用中间件2……->调用真正的业务处理逻辑
	router.Handler("GET", "/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(Index)))
	//handle里的panic捕获输出
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error:%s", v)
	}

	playRouter := httprouter.New()
	playRouter.GET("/", UserInfo)

	toolRouter := httprouter.New()
	toolRouter.GET("/", UserInfo)

	hs := make(HostSwitch)
	hs["play.flysnow.org:12345"] = playRouter
	hs["tool.flysnow.org:12345"] = toolRouter

	log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(http.ListenAndServe(":12345", hs))
}
