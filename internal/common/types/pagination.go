package types

type PaginatedResponse[T any] struct {
	Data any                `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}

type PaginationMetadata struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
}
