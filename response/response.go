package response

import (
	"myslb/constant"
	"net/http"
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

//第一个返回结构体
type CGIResponse struct {
	Content string `json:"content"`
}

type HealthCheckResponse struct {
	Status int `json:"status"`
}

func respHandler(code int, msg string, results ...interface{}) (r *Response) {
	r = &Response{
		Code:   code,
		Msg:    msg,
	}
	if len(results) > 0 {
		r.Result = results[0]
	}
	return r
}

func ErrorResp(code int, msg string) *Response {
	return respHandler(code, msg)
}

func SuccessResp(result ...interface{}) *Response {
	return respHandler(http.StatusOK, constant.MsgOK, result...)
}