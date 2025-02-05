package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-gin-demo/docs"
	"go-gin-demo/global"
	"go-gin-demo/middleware"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type IfnRegistRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var (
	gfnRoutes []IfnRegistRoute
)

func RegistRoute(fn IfnRegistRoute) {
	if fn == nil {
		return
	}

	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {

	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancelCtx()

	r := gin.Default()
	//跨域
	r.Use(middleware.CorsMiddleware())
	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")

	//鉴权
	//rgAuth.Use(middleware.Auth())

	initBasePlatformRouters()

	customValidator()

	for _, fnRegistRoute := range gfnRoutes {
		fnRegistRoute(rgPublic, rgAuth)
	}

	//集成swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	stPort := viper.GetString("server.port")
	if stPort == "" {
		stPort = "8999"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", stPort),
		Handler: r,
	}

	go func() {
		global.Logger.Info("start Listen: ", stPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error("start server Error: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()

	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error("stop server" + err.Error())
		return
	}

	global.Logger.Info("stop server success")
}

func initBasePlatformRouters() {
	//InitUserRouters()
	//InitUpdateVersionRouters()
	//InitMqRouters()
	//InitItemRouters()
	//InitSSHRouters()
	InitCaseRouters()
	//InitWsRouters()
}

func customValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("first_is_a", func(fl validator.FieldLevel) bool {
			if value, ok := fl.Field().Interface().(string); ok {
				if value != "" && 0 == strings.Index(value, "a") {
					return true
				}
			}
			return false
		})
	}
}
