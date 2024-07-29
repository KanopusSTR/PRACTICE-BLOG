package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"server/internal/models"
	"strconv"
)

func (h *handler) DeletePost(c *gin.Context) {
	fun := func() (models.DeletePost, error) {
		dp := models.DeletePost{}
		id, err := strconv.Atoi(c.Param("id"))
		mail := c.GetString("mail")
		dp.PostId = id
		dp.Mail = mail
		return dp, err
	}
	c.JSON(h.handlerS.DeletePost(fun))
}

func (h *handler) EditPost(c *gin.Context) {
	fun := func() (models.EditPost, error) {
		ep := models.EditPost{}
		err := c.BindJSON(&ep)
		if err != nil {
			return ep, err
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return ep, errors.WithMessage(err, "invalid post id")
		}
		ep.PostId = id
		ep.Mail = c.GetString("mail")
		return ep, nil
	}
	c.JSON(h.handlerS.EditPost(fun))
}

func (h *handler) GetPost(c *gin.Context) {
	fun := func() (models.GetPost, error) {
		gp := models.GetPost{}
		postId, err := strconv.Atoi(c.Param("id"))
		gp.Id = postId
		return gp, err
	}
	c.JSON(h.handlerS.GetPost(fun))
}

func (h *handler) GetPosts(c *gin.Context) {
	c.JSON(h.handlerS.GetPosts())
}

func (h *handler) WritePost(c *gin.Context) {
	fun := func() (models.WritePost, error) {
		wp := models.WritePost{}
		err := c.BindJSON(&wp)
		wp.Mail = c.GetString("mail")
		return wp, err
	}
	c.JSON(h.handlerS.WritePost(fun))
}
