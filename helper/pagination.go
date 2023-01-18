package helper

import "golang/golang-skeleton/dto/pagination"

func GeneratePaginationFromRequest(paginationRequest *pagination.PaginationRequest) pagination.Pagination {

	page := 1
	limit := 15
	sort := "created_at asc"

	if paginationRequest.Page != 0 {
		page = paginationRequest.Page
	}
	if paginationRequest.Per_Page != 0 {
		limit = paginationRequest.Per_Page
	}
	if paginationRequest.OrderBy != "" {
		sort = paginationRequest.OrderBy
	}

	return pagination.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func BuildMetadataResponse(total int64, per_page int, page int) pagination.Metadata {
	var metadata pagination.Metadata

	metadata.Per_Page = per_page
	metadata.Page = page
	metadata.Total = total

	return metadata
}
