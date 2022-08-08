package types

type PaginateMetadata struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	TotalItems  int `json:"totalItems"`
}

type IGetOverviewQuery struct {
	Page  int      `query:"page"`
	Limit int      `query:"limit"`
	Id    []string `query:"id"`
}
