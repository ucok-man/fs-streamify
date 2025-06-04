package models

import "math"

type Metadata struct {
	CurrentPage  int64 `json:"current_page"`
	PageSize     int64 `json:"page_size"`
	FirstPage    int64 `json:"first_page"`
	LastPage     int64 `json:"last_page"`
	TotalRecords int64 `json:"total_records"`
}

func calculateMetadata(totalRecords, page, pageSize int64) Metadata {
	if totalRecords == 0 {
		return Metadata{} // return an empty Metadata struct if there are no records
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int64(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
