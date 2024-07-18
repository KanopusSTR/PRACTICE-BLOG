package handlers

import (
	"github.com/gin-gonic/gin"
	"server/internal/models"
)

func (h *handler) GetUser(c *gin.Context) {
	fun := func() (models.GetUser, error) {
		gu := models.GetUser{}
		gu.Mail = c.Param("mail")
		return gu, nil
	}
	c.JSON(h.app.Handler.GetUser(fun))
}
