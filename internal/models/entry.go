package models

type RegisterRequest struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type LoginMiddleware struct {
	Token string `json:"token"`
}
