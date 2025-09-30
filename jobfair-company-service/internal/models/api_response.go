package models

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}
