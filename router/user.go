package router

import "gin_chat/service"

func UserRoutes(routerBase *RouterApp) {

	// api 路由
	routerBase.ApiGroup.GET("/index", service.GetIndexView)
	routerBase.ApiGroup.POST("/register", service.Register)
	routerBase.ApiGroup.POST("/login", service.Login)
}
