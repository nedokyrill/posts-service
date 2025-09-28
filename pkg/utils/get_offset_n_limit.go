package utils

import "github.com/nedokyrill/posts-service/pkg/consts"

func GetOffsetNLimit(page *int32, pageSize int) (int, int) {
	var newPage int32

	if page == nil || *page <= 0 {
		newPage = 0
	} else {
		newPage = *page - 1
	}

	if pageSize <= 0 {
		pageSize = consts.PageSize
	}

	offset := int(newPage) * pageSize
	limit := pageSize

	return offset, limit
}
