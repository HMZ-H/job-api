package models

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
	Errors  []string    `json:"errors,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Object     interface{} `json:"object"`
	PageNumber int         `json:"page_number"`
	PageSize   int         `json:"page_size"`
	TotalSize  int64       `json:"total_size"`
	Errors     []string    `json:"errors,omitempty"`
}
