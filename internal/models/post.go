package models

type PostRequest struct {
	Header string `json:"header"`
	Body   string `json:"body"`
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
