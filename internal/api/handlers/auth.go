package handlers

import (
	"github.com/gin-gonic/gin"
	"server/internal/models"
)

func (h *handler) Login(c *gin.Context) {
	fun := func() (models.LoginRequest, error) {
		var req models.LoginRequest
		return req, c.BindJSON(&req)
	}
	c.JSON(h.app.Handler.Login(fun))
}

func (h *handler) Register(c *gin.Context) {
	fun := func() (models.RegisterRequest, error) {
		var req models.RegisterRequest
		return req, c.BindJSON(&req)
	}
	c.JSON(h.app.Handler.Register(fun))
}
