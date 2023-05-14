package route

import (
	"github.com/gin-gonic/gin"
	"go-prerender/server/controller"
	"go-prerender/server/route/middleware/cors"
)

func SetupRouter(engine *gin.Engine) {
	//设置路由中间件
	engine.Use(cors.SetUp()) // 跨域

	engine.GET("/", controller.Index)            // 首页
	engine.GET("/proxy/*name", controller.Proxy) // 代理

	// 全部反向代理
	engine.NoRoute(controller.Proxy)

	engine.NoMethod(func(c *gin.Context) {

	})

}
