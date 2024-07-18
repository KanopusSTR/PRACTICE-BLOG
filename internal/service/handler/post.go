package handlers

import (
	"net/http"
	"server/internal/models"
	"time"
)

func (hs *handlerService) DeletePost(fun func() (models.DeletePost, error)) (int, models.ResultResponseBody) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.ResultResponseBody{Message: "postId is required and it must be non-negative"}
	}
	if post, err := hs.users.GetPost(req.PostId); err != nil {
		return http.StatusNotFound, models.ResultResponseBody{Message: "deletePost error: " + err.Error()}
	} else if post.AuthorMail != req.Mail {
		return http.StatusForbidden, models.ResultResponseBody{Message: "deletePost error: you do not have permission"}
	}
	if err := hs.users.DeletePost(req.PostId); err != nil {
		return http.StatusBadRequest, models.ResultResponseBody{Message: "deletePost error: " + err.Error()}
	}
	return http.StatusOK, models.ResultResponseBody{Message: "success"}
}

func (hs *handlerService) EditPost(fun func() (models.EditPost, error)) (int, models.ResultResponseBody) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.ResultResponseBody{Message: "editPost error: " + err.Error()}
	}

	if post, err := hs.users.GetPost(req.PostId); err != nil {
		return http.StatusNotFound, models.ResultResponseBody{Message: "editPost error: " + err.Error()}
	} else if post.AuthorMail != req.Mail {
		return http.StatusForbidden, models.ResultResponseBody{Message: "editPost error: you do not have permission"}
	}
	if err := hs.users.EditPost(req.PostId, &req.Header, &req.Body); err != nil {
		return http.StatusNotFound, models.ResultResponseBody{Message: "editPost error: " + err.Error()}
	}
	return http.StatusOK, models.ResultResponseBody{Message: "success"}
}

func (hs *handlerService) GetPost(fun func() (models.GetPost, error)) (int, models.GetPostResponse) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.GetPostResponse{Message: "postId is required and it must be non-negative"}
	}
	post, err := hs.users.GetPost(req.Id)
	if err != nil {
		return http.StatusNotFound, models.GetPostResponse{
			Message: "getPost error: " + err.Error()}
	}
	return http.StatusOK, models.GetPostResponse{Message: "success", Posts: *post}
}

func (hs *handlerService) GetPosts() (int, models.GetPostsResponse) {
	posts := hs.users.GetPosts()
	return http.StatusOK, models.GetPostsResponse{Posts: posts}
}

func (hs *handlerService) WritePost(fun func() (models.WritePost, error)) (int, models.WritePostResponse) {
	req, err := fun()
	if err != nil {
		return http.StatusBadRequest, models.WritePostResponse{Message: "writePost error: " + err.Error()}
	}
	postId, err := hs.users.WritePost(&req.Header, &req.Body, time.Now(), req.Mail)
	if err != nil {
		return http.StatusBadRequest, models.WritePostResponse{Message: "writePost error: " + err.Error()}
	}
	return http.StatusOK, models.WritePostResponse{Message: "success", PostId: postId}
}
