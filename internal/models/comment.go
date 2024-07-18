package models

type WriteCommentRequest struct {
	Text   string `json:"text"`
	PostId int    `json:"post_id"`
}

type WriteComment struct {
	Text   string `json:"text"`
	PostId int
	Mail   string
}

type DeleteComment struct {
	PostId    int
	CommentId int
	Mail      string
}

type GetCommentsRequest struct {
	PostId int `json:"post_id"`
}

type GetCommentsResponse struct {
	Message  string        `json:"message"`
	Comments []interface{} `json:"comments"`
}
