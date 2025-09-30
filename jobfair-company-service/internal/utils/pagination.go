package utils

import "math"

type PaginationParams struct {
	Page  int
	Limit int
}

func NewPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func CalculateTotalPages(totalItems int64, limit int) int {
	return int(math.Ceil(float64(totalItems) / float64(limit)))
}