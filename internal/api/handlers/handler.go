package handlers

import (
	"github.com/gin-gonic/gin"
	"server/internal/service"
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
	app *service.I
}

func New(i *service.I) Handler {
	return &handler{i}
}

func (h *handler) HandleSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
