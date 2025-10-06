package types

type QueryParams struct {
	Page    int    `binding:"required,number,min=1" json:"page" form:"page"`
	Limit   int    `binding:"required,number,min=1,max=100" json:"limit" form:"limit"`
	Search  string `json:"search" form:"search"`
	SortBy  string `json:"sortBy" form:"sortBy"`
	OrderBy string `binding:"omitempty,oneof=ASC DESC" json:"orderBy" form:"orderBy"`
}
