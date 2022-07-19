package requester

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"myslb/response"
	"net/http"
	"reflect"
)

type Requester interface {
	Request() (resp *response.Response)
}

func Handler(r Requester, c *gin.Context) {
	if reflect.TypeOf(r).Kind() != reflect.Ptr {
		fmt.Println("Handler request ptr error")
		return
	}
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := c.ShouldBindWith(r, b); err != nil {
		fmt.Println("Parse request error")
		fmt.Println(err)
	}
	result := r.Request()
	c.JSON(http.StatusOK, result)
}

// 实现一个Requester
type GetCGIHandler struct {}

func (r *GetCGIHandler) Request() *response.Response {
	resp := &response.CGIResponse{Content: "GOCGI"}
	return response.SuccessResp(resp)
}