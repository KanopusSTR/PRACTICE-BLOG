package models

import "server/internal/entities"

type PostRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

type GetPostsResponse struct {
	Posts []interface{} `json:"posts"`
}

type GetPostResponse struct {
	Message string        `json:"message"`
	Posts   entities.Post `json:"post"`
}

type WritePostResponse struct {
	Message string `json:"message"`
	PostId  int    `json:"post_id"`
}

type DeletePost struct {
	PostId int    `json:"post_id"`
	Mail   string `json:"mail"`
}

type EditPost struct {
	Header string `json:"header"`
	Body   string `json:"body"`
	PostId int    `json:"post_id"`
	Mail   string `json:"mail"`
}

type GetPost struct {
	Id int `json:"id"`
}

type WritePost struct {
	Header string `json:"header"`
	Body   string `json:"body"`
	Mail   string `json:"mail"`
}
