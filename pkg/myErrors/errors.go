package myErrors

import "errors"

var (
	EmptyPost         = errors.New("header or body of post is empty")
	PostNotFound      = errors.New("post not found")
	IncorrectPassword = errors.New("incorrect password")
	CommentNotFound   = errors.New("comment not found")
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user already exists")
	InvalidMail       = errors.New("mail is invalid")
	EmptyRegisterData = errors.New("username or password is empty")
	EmptyField        = errors.New("one or more fields are empty")
)
