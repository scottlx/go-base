package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var cookieMap map[string]*http.Cookie

func updateCookieRouting() {
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			for vrouterId := range cookieMap {
				cookieMap[vrouterId] = login("scott", "lin")
			}
		}
	}()
}
func login(username, password string) *http.Cookie {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.PostForm("http://localhost:9292/user/login", url.Values{
		"username": {username},
		"password": {password},
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("登录失败，请检查用户名密码")
		return nil
	}

	fmt.Println(resp.StatusCode)
	fmt.Printf("登录返回cookie： %s \n", resp.Cookies())
	fmt.Printf("登录返回header： %v \n", resp.Header)

	if len(resp.Cookies()) == 0 {
		return nil
	}
	return resp.Cookies()[0]
}

func proxy(c *gin.Context) {
	remote, err := url.Parse("http://localhost:9292")
	if err != nil {
		panic(err)
	}
	if cookieMap[c.Param("vrouterId")] == nil {
		cookieMap[c.Param("vrouterId")] = login("scott", "lin")
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Header.Set("Referer", req.URL.String()) // 供vrouter web 重定向使用，重定向到被代理的地址

		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
		req.AddCookie(cookieMap[c.Param("vrouterId")])

		fmt.Println(req.Header)
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {

	r := gin.Default()
	cookieMap = make(map[string]*http.Cookie)
	//testCookie()
	updateCookieRouting()

	r.Any("/sso/:vrouterId/*proxyPath", proxy)

	r.Run(":8080")
}
