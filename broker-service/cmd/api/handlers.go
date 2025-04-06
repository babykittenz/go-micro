package main

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}
