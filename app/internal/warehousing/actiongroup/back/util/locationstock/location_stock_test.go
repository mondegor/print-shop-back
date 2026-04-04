package locationstock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/dto"
	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/util/locationstock"
)

func Test_ChunkByLocationID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    []dto.LocationStock
		want [][]dto.LocationStock
	}{
		{
			name: "test1",
			s:    nil,
			want: [][]dto.LocationStock{},
		},
		{
			name: "test2",
			s: []dto.LocationStock{
				{
					LocationID:  1,
					ContainerID: 1,
				},
			},
			want: [][]dto.LocationStock{
				{
					{
						LocationID:  1,
						ContainerID: 1,
					},
				},
			},
		},
		{
			name: "test3",
			s: []dto.LocationStock{
				{
					LocationID:  1,
					ContainerID: 1,
				},
				{
					LocationID:  1,
					ContainerID: 2,
				},
			},
			want: [][]dto.LocationStock{
				{
					{
						LocationID:  1,
						ContainerID: 1,
					},
					{
						LocationID:  1,
						ContainerID: 2,
					},
				},
			},
		},
		{
			name: "test4",
			s: []dto.LocationStock{
				{
					LocationID:  1,
					ContainerID: 1,
				},
				{
					LocationID:  2,
					ContainerID: 2,
				},
			},
			want: [][]dto.LocationStock{
				{
					{
						LocationID:  1,
						ContainerID: 1,
					},
				},
				{
					{
						LocationID:  2,
						ContainerID: 2,
					},
				},
			},
		},
		{
			name: "test5",
			s: []dto.LocationStock{
				{
					LocationID:  1,
					ContainerID: 1,
				},
				{
					LocationID:  1,
					ContainerID: 2,
				},
				{
					LocationID:  2,
					ContainerID: 2,
				},
			},
			want: [][]dto.LocationStock{
				{
					{
						LocationID:  1,
						ContainerID: 1,
					},
					{
						LocationID:  1,
						ContainerID: 2,
					},
				},
				{
					{
						LocationID:  2,
						ContainerID: 2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ss := make([][]dto.LocationStock, 0)
			for s := range locationstock.ChunkByLocationID(tt.s) {
				ss = append(ss, s)
			}

			assert.Equal(t, tt.want, ss)
		})
	}
}
