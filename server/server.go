package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go-prerender/internal"
	"go-prerender/internal/config"
	"go-prerender/server/route"
)

func StartHttpService() {
	gin.SetMode("release")
	ginEngine := gin.New()

	// 性能分析 - 正式环境不要使用！！！

	pprof.Register(ginEngine, "dev/pprof")
	//ginEngine.Use(ginprom.PromMiddleware(nil))

	// register the `/metrics` route.
	//ginEngine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	// 设置路由
	route.SetupRouter(ginEngine)

	webApiPort := config.GetConfig().HttpServ.Port
	internal.SLogger.GetStdoutLogger().Infof("the api service has been started. addr: %s%s", "http://127.0.0.1:", webApiPort)

	var err error
	err = ginEngine.Run(":" + webApiPort)
	if err != nil {
		internal.SLogger.GetStdoutLogger().Fatalf("the api service failed to start! %s", err.Error())
	}
}
