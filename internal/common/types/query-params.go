package types

type QueryParams struct {
	Search     string `form:"search"`
	SortDir    string `form:"sortDir"`
	SortBy     string `form:"sortBy"`
	IncludeIds *[]int `form:"includeIds[]" default:"[]"` // Optional Field to include specific IDs in the query results
	ExcludeIds *[]int `form:"excludeIds[]" default:"[]"` // Optional Field to exclude specific IDs from the query results
	Limit      int    `form:"limit" binding:"min=1" `
}

type PaginationQueryParams struct {
	QueryParams
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"pageSize" binding:"min=1,max=100"`
}

func NewPaginationQueryParams() PaginationQueryParams {
	return PaginationQueryParams{
		QueryParams: QueryParams{
			Search:     "",
			SortDir:    "DESC",
			SortBy:     "created_at",
			IncludeIds: &[]int{},
			ExcludeIds: &[]int{},
			Limit:      10,
		},
		Page:     1,
		PageSize: 10,
	}
}
