package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/models"
)

func (h *handler) LoginMiddleware(c *gin.Context) {
	fun := func() (models.LoginMiddleware, error) {
		t := models.LoginMiddleware{}
		token := c.Request.Header.Get("JWT-Token")
		t.Token = token
		return t, nil
	}
	code, rb, mail := h.app.Handler.LoginMiddleware(fun)
	if code != http.StatusOK {
		c.JSON(code, rb)
		c.Abort()
		return
	}
	c.Set("mail", mail)
	c.Next()
}
