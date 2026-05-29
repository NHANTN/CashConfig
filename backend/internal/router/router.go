package router

import (
	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/handler"
	"github.com/cashier-config/server/internal/middleware"
	"github.com/cashier-config/server/internal/service"
)

func Setup(handlers []handler.Handler, authSvc *service.AuthService, apiKey string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	authed := api.Group("")
	authed.Use(middleware.JWTAuth(authSvc, apiKey))

	for _, h := range handlers {
		h.Register(api, authed)
	}

	return r
}
