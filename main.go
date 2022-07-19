package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myslb/requester"
)

func main() {
	r := InitRouters()
	if err := r.Run(":33802"); err != nil {
		fmt.Println("run gin app error")
		fmt.Println(err)
	}
}

func InitRouters() *gin.Engine {
	router := gin.New()
	ver := fmt.Sprintf("/%s", "v1")
	routerGroupVer := router.Group(ver)

	routerGroupVer.GET("/cgi", Handler(&requester.GetCGIHandler{}))
	return router
}

func Handler(r requester.Requester) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester.Handler(r, c)
	}
}