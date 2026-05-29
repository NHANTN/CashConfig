package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/cashier-config/server/internal/handler"
	"github.com/cashier-config/server/internal/middleware"
	"github.com/cashier-config/server/internal/service"
)

func Setup(handlers []handler.Handler, authSvc *service.AuthService, apiKey string) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	authed := api.Group("")
	authed.Use(middleware.JWTAuth(authSvc, apiKey))

	for _, h := range handlers {
		h.Register(api, authed)
	}

	return r
}
