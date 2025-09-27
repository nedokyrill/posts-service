package utils

import "github.com/nedokyrill/posts-service/pkg/consts"

func GetOffsetNLimit(page *int, pageSize int) (int, int) {
	var newPage int

	if page == nil || *page <= 0 {
		newPage = 1
	} else {
		newPage = *page - 1
	}

	if pageSize <= 0 {
		pageSize = consts.PageSize
	}

	offset := newPage * pageSize
	limit := pageSize

	return offset, limit
}
