package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	n := negroni.Classic()

	router := mux.NewRouter()
	router.Handle("/", handler())
	router.HandleFunc("/flysnow", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "Blog:www.flysnow.org\n")
		io.WriteString(rw, "Wechat:flysnow_org")
	})
	n.UseHandler(router)
	n.Run(":1234")
}

func handler() http.Handler {
	return http.HandlerFunc(myHandler)
}

func myHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	io.WriteString(rw, "Hello World")
}
