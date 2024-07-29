package handlers

import (
	"github.com/gin-gonic/gin"
	handlers "server/internal/service/handler"
)

type Handler interface {
	WritePost(c *gin.Context)
	GetPost(c *gin.Context)
	GetPosts(c *gin.Context)
	EditPost(c *gin.Context)
	DeletePost(c *gin.Context)

	GetUser(ctx *gin.Context)

	LoginMiddleware(ctx *gin.Context)

	WriteComment(c *gin.Context)
	GetComments(c *gin.Context)
	DeleteComment(c *gin.Context)

	Login(c *gin.Context)
	Register(c *gin.Context)
}

type handler struct {
	handlerS handlers.Service
}

func New(h handlers.Service) Handler {
	return &handler{handlerS: h}
}

func (h *handler) HandleSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
