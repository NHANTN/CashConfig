package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(api gin.IRouter, authed gin.IRouter)
}
