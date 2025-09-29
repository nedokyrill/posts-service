package utils

import (
	"testing"

	"github.com/nedokyrill/posts-service/pkg/consts"
)

func ptr(i int32) *int32 { return &i }

func TestGetOffsetNLimit(t *testing.T) {
	tests := []struct {
		name     string
		page     *int32
		pageSize int
		wantOff  int
		wantLim  int
	}{
		{
			name:     "nil page, zero pageSize",
			page:     nil,
			pageSize: 0,
			wantOff:  0,
			wantLim:  consts.PageSize,
		},
		{
			name:     "page 1, pageSize 10",
			page:     ptr(1),
			pageSize: 10,
			wantOff:  0,
			wantLim:  10,
		},
		{
			name:     "page 2, pageSize 10",
			page:     ptr(2),
			pageSize: 10,
			wantOff:  10,
			wantLim:  10,
		},
		{
			name:     "page 3, pageSize 5",
			page:     ptr(3),
			pageSize: 5,
			wantOff:  10,
			wantLim:  5,
		},
		{
			name:     "negative page",
			page:     ptr(-1),
			pageSize: 10,
			wantOff:  0,
			wantLim:  10,
		},
		{
			name:     "page 2, zero pageSize",
			page:     ptr(2),
			pageSize: 0,
			wantOff:  1 * consts.PageSize,
			wantLim:  consts.PageSize,
		},
		{
			name:     "nil page, negative pageSize",
			page:     nil,
			pageSize: -5,
			wantOff:  0,
			wantLim:  consts.PageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOff, gotLim := GetOffsetNLimit(tt.page, tt.pageSize)
			if gotOff != tt.wantOff || gotLim != tt.wantLim {
				t.Fatalf("GetOffsetNLimit(%v, %d) = (offset=%d, limit=%d), want (offset=%d, limit=%d)",
					tt.page, tt.pageSize, gotOff, gotLim, tt.wantOff, tt.wantLim)
			}
		})
	}
}
