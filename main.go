package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myslb/loadbalance"
	"myslb/requester"
	"net/http"
	"net/http/httputil"
)

func main() {
	r := InitRouters()
	if err := r.Run(":33803"); err != nil {
		fmt.Println("run gin app error")
		fmt.Println(err)
	}
}

func InitRouters() *gin.Engine {
	router := gin.New()
	ver := fmt.Sprintf("/%s", "v1")
	routerGroupVer := router.Group(ver)

	// init loadbalance
	//rr := &loadbalance.RoundRobin{ServerPool: loadbalance.NewRRPool()}
	//fmt.Printf("%+v\n", rr.ServerPool)

	routerGroupVer.GET("/cgi", Handler(&requester.GetCGIHandler{}))
	routerGroupVer.GET("/healthCheck", Handler(&requester.HealthCheckHandler{}))
	//routerGroupVer.GET("/:name", ReverseProxy(rr))

	return router
}

func Handler(r requester.Requester) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester.Handler(r, c)
	}
}

func ReverseProxy(rr *loadbalance.RoundRobin) gin.HandlerFunc {
	return func(c *gin.Context) {
		server, _ := rr.GetNextPeer()
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = server
			req.Host = server
		}
		fmt.Printf("Proxy %s\n", server)
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}