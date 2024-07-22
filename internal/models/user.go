package models

type GetUser struct {
	Mail string `json:"mail"`
}

type ProfileResponse struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}
