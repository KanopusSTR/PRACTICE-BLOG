package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"server/internal/models"
	"strconv"
)

func (h *handler) DeleteComment(c *gin.Context) {
	fun := func() (models.DeleteComment, error) {
		var req models.DeleteComment
		err := c.BindJSON(&req)
		if err != nil {
			return req, err
		}
		commentId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return req, errors.WithMessage(err, "invalid comment id")
		}
		req.CommentId = commentId
		req.Mail = c.GetString("mail")
		return req, nil
	}
	c.JSON(h.app.Handler.DeleteComment(fun))
}

func (h *handler) GetComments(c *gin.Context) {
	fun := func() (models.GetCommentsRequest, error) {
		var req models.GetCommentsRequest
		return req, c.BindJSON(&req)
	}
	c.JSON(h.app.Handler.GetComments(fun))
}

func (h *handler) WriteComment(c *gin.Context) {
	fun := func() (models.WriteComment, error) {
		var req models.WriteComment
		err := c.BindJSON(&req)
		if err != nil {
			return req, err
		}
		req.Mail = c.GetString("mail")
		return req, nil
	}
	c.JSON(h.app.Handler.WriteComment(fun))
}
