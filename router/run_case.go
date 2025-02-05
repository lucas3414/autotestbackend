package router

import (
	"github.com/gin-gonic/gin"
	"go-gin-demo/api"
)

func InitCaseRouters() {
	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {

		runCaseApi := api.NewRunCaseApi()
		rgPublicRunCase := rgPublic.Group("AutoUiTest")
		{
			rgPublicRunCase.POST("/runCase", runCaseApi.RunCase)

		}

	})
}
