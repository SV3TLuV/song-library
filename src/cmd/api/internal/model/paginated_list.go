package model

type PaginatedList[T any] struct {
	Page       uint
	PageSize   uint
	TotalPages uint
	Items      []T
}
