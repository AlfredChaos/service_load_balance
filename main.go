package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myslb/requester"
	"net/http"
	"net/http/httputil"
)

func main() {
	r := InitRouters()
	if err := r.Run(":3380"); err != nil {
		fmt.Println("run gin app error")
		fmt.Println(err)
	}
}

func InitRouters() *gin.Engine {
	router := gin.New()
	ver := fmt.Sprintf("/%s", "v1")
	routerGroupVer := router.Group(ver)

	//routerGroupVer.GET("/cgi", Handler(&requester.GetCGIHandler{}))
	routerGroupVer.GET("/:name", ReverseProxy())
	return router
}

func Handler(r requester.Requester) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester.Handler(r, c)
	}
}

func ReverseProxy() gin.HandlerFunc {
	target := "127.0.0.1:33802"
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Host = target
		}
		fmt.Println("Proxy here")
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}