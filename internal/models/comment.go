package models

type WriteCommentRequest struct {
	Text   string `json:"text"`
	PostId int    `json:"post_id"`
}

type WriteComment struct {
	Text   string `json:"text"`
	PostId int    `json:"post_id"`
	Mail   string `json:"mail"`
}

type DeleteComment struct {
	PostId    int    `json:"post_id"`
	CommentId int    `json:"comment_id"`
	Mail      string `json:"mail"`
}

type GetCommentsRequest struct {
	PostId int `json:"post_id"`
}
