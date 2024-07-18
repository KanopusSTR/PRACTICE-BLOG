package models

type GetProfileResponse struct {
	Message string          `json:"message"`
	Profile ProfileResponse `json:"profile"`
}

type GetUser struct {
	Mail string `json:"mail"`
}

type ProfileResponse struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}
