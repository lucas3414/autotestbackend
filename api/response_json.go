package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ResponseJson struct {
	Status int    `json:"-"`
	Code   int    `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Data   any    `json:"data,omitempty"`
	Total  int    `json:"total,omitempty"`
}

func (m ResponseJson) IsEmpty() bool {
	return reflect.DeepEqual(m, ResponseJson{})
}

func HttpResponse(ctx *gin.Context, status int, resp ResponseJson) {
	if resp.IsEmpty() {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.AbortWithStatusJSON(status, resp)
}

func buildStatus(resp ResponseJson, defaultStatus int) int {
	if resp.Status == 0 {
		return defaultStatus
	}
	return resp.Status
}

func OK(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusBadRequest), resp)
}

func ServerFail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusInternalServerError), resp)
}
