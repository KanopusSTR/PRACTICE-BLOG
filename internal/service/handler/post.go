package handlers

import (
	"net/http"
	"server/internal/models"
	"time"
)

func (hs *handlerService) DeletePost(fun func() (models.DeletePost, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "postId is required and it must be non-negative"}
	}
	if post, err := hs.users.GetPost(req.PostId); err != nil {
		return http.StatusNotFound, models.Response{Message: "deletePost error: " + err.Error()}
	} else if post.AuthorMail != req.Mail {
		return http.StatusForbidden, models.Response{Message: "deletePost error: you do not have permission"}
	}
	if err := hs.users.DeletePost(req.PostId); err != nil {
		return http.StatusBadRequest, models.Response{Message: "deletePost error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success"}
}

func (hs *handlerService) EditPost(fun func() (models.EditPost, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "editPost error: " + err.Error()}
	}

	if post, err := hs.users.GetPost(req.PostId); err != nil {
		return http.StatusNotFound, models.Response{Message: "editPost error: " + err.Error()}
	} else if post.AuthorMail != req.Mail {
		return http.StatusForbidden, models.Response{Message: "editPost error: you do not have permission"}
	}
	if err := hs.users.EditPost(req.PostId, &req.Header, &req.Body); err != nil {
		return http.StatusNotFound, models.Response{Message: "editPost error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success"}
}

func (hs *handlerService) GetPost(fun func() (models.GetPost, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "postId is required and it must be non-negative"}
	}
	post, err := hs.users.GetPost(req.Id)
	if err != nil {
		return http.StatusNotFound, models.Response{
			Message: "getPost error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success", Data: *post}
}

func (hs *handlerService) GetPosts() (int, models.Response) {
	posts := hs.users.GetPosts()
	return http.StatusOK, models.Response{Data: posts}
}

func (hs *handlerService) WritePost(fun func() (models.WritePost, error)) (int, models.Response) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "writePost error: " + err.Error()}
	}
	postId, err := hs.users.WritePost(&req.Header, &req.Body, time.Now(), req.Mail)
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "writePost error: " + err.Error()}
	}
	return http.StatusOK, models.Response{Message: "success", Data: postId}
}
