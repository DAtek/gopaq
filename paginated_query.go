package gopaq

import (
	"gorm.io/gorm"
)

var DefaultLimit = 50

type PaginatedQueryResult[T interface{}] struct {
	Total uint
	Items T
}

func FindWithPagination[T interface{}](
	query *gorm.DB,
	items T,
	page uint,
	pageSize uint,
) (*PaginatedQueryResult[T], error) {
	limit := int(pageSize)
	if limit == 0 {
		limit = DefaultLimit
	}

	offset := (int(page) - 1)
	if offset == -1 {
		offset = 0
	}

	offset = offset * int(pageSize)

	total := int64(0)
	if r := query.Count(&total); r.Error != nil {
		return nil, r.Error
	}

	result := &PaginatedQueryResult[T]{
		Items: items,
		Total: uint(total),
	}

	r := query.Limit(int(limit)).Offset(int(offset)).Find(&result.Items)
	return result, r.Error
}
