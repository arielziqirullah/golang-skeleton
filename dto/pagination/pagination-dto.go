package pagination

type DataResponse struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
}

type Metadata struct {
	Per_Page int   `json:"per_page"`
	Page     int   `json:"page"`
	Total    int64 `json:"total"`
}

type PaginationRequest struct {
	Per_Page int    `json:"per_page" form:"per_page"`
	Page     int    `json:"page" form:"page"`
	OrderBy  string `json:"order_by" form:"order_by"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
