package router

import (
	"gin_chat/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// 匿名导入生成的接口文档包
	_ "gin_chat/docs"
)

type RouterApp struct {
	ApiGroup *gin.RouterGroup
	router   *gin.Engine
}

func InitRouter() *gin.Engine {
	systemSettings := viper.GetStringMapString("system")
	gin.SetMode(systemSettings["env"])
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware())
	routerGroup := router.Group("api") // 地址前缀
	routerGroupApp := RouterApp{
		router:   router,
		ApiGroup: routerGroup,
	}
	routerGroupApp.registerRouter()

	return router
}

func (routerBase *RouterApp) registerRouter() {
	// 注册swagger静态文件路由
	routerBase.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	UserRoutes(routerBase)
}
