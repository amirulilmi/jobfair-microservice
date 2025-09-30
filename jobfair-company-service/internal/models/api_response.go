package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, code string, details interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func PaginatedSuccessResponse(message string, data interface{}, meta PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: meta,
	}
}
