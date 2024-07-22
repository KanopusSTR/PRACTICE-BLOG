package models

type ResultResponseBody struct {
	Message string `json:"message"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
