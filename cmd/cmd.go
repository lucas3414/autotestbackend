package cmd

import (
	"fmt"
	"go-gin-demo/conf"
	"go-gin-demo/global"
	"go-gin-demo/router"
)

func Start() {

	var initErr error

	conf.InitConfig()
	global.Logger = conf.InitLogger()

	// 初始化过程中遇到问题的处理
	if initErr != nil {
		if global.Logger != nil {
			global.Logger.Error(initErr.Error())
		}
		panic(initErr.Error())
	}

	router.InitRouter()
}

func Clean() {
	fmt.Println("=======Clean=========")
}
