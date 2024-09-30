package imposition_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_mrlog "github.com/mondegor/go-webcore/mrlog/mock"
	"github.com/stretchr/testify/assert"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

// TODO: ДОБАВИТЬ ТЕСТ MIRROR

func newLoggerStub(ctrl *gomock.Controller) *mock_mrlog.MockLogger {
	mockLogger := mock_mrlog.NewMockLogger(ctrl)
	mockLoggerEvent := mock_mrlog.NewMockLoggerEvent(ctrl)

	mockLogger.EXPECT().Debug().Return(mockLoggerEvent).MinTimes(1)
	mockLoggerEvent.EXPECT().MsgFunc(gomock.Any()).MinTimes(1)

	return mockLogger
}

func TestImposition_CalcWithAllowRotation(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	tests := []struct {
		name string
		item rect.Item
		out  rect.Format
		want imposition.AlgoResult
	}{
		// WARNING: большинство данных совпадает с данными теста CalcTotal
		{
			name: "test1",
			item: rect.Item{Format: rect.Format{Width: 88, Height: 48}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 435, Height: 300}, // 450x320 - frame 15x20
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 408, Height: 298},
				Fragments: base.Fragments{{ByWidth: 4, ByHeight: 6}, {ByWidth: 3, ByHeight: 1}},
				Total:     27,
				RestArea:  10356,
			},
		},
		{
			name: "test2",
			item: rect.Item{Format: rect.Format{Width: 88, Height: 48}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 435, Height: 290}, // 450x320 - frame 15x30
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 408, Height: 268},
				Fragments: base.Fragments{{ByWidth: 4, ByHeight: 5}, {ByWidth: 3, ByHeight: 1}},
				Total:     23,
				RestArea:  23966,
			},
		},
		{
			name: "test3",
			item: rect.Item{Format: rect.Format{Width: 48, Height: 88}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 435, Height: 300}, // 450x320 - frame 15x20
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 398, Height: 268},
				Fragments: base.Fragments{{ByWidth: 8, ByHeight: 3}},
				Total:     24,
				RestArea:  23836,
			},
		},
		{
			name: "test4",
			item: rect.Item{Format: rect.Format{Width: 48, Height: 88}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 435, Height: 290}, // 450x320 - frame 15x30
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 398, Height: 268},
				Fragments: base.Fragments{{ByWidth: 8, ByHeight: 3}},
				Total:     24,
				RestArea:  19486,
			},
		},
		{
			name: "test5",
			item: rect.Item{Format: rect.Format{Width: 420, Height: 300}},
			out:  rect.Format{Width: 435, Height: 300}, // 450x320 - frame 15x20
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 420, Height: 300},
				Fragments: base.Fragments{{ByWidth: 1, ByHeight: 1}},
				Total:     1,
				RestArea:  4500,
			},
		},
		{
			name: "test6",
			item: rect.Item{Format: rect.Format{Width: 210, Height: 297}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 435, Height: 300}, // 450x320 - frame 15x20
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 422, Height: 297},
				Fragments: base.Fragments{{ByWidth: 2, ByHeight: 1}},
				Total:     2,
				RestArea:  5166,
			},
		},
		{
			// данный тест не зависит от дистанции между изделиями, т.к. изделие всего одно
			name: "test7",
			item: rect.Item{Format: rect.Format{Width: 235, Height: 100}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 235, Height: 100}, // 450x320 - frame 15x20
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 235, Height: 100},
				Fragments: base.Fragments{{ByWidth: 1, ByHeight: 1}},
				Total:     1,
				RestArea:  0,
			},
		},
		{
			// height = y = 2x + 1; width = x + y + 3; при x = 100 -> y = 201, height = 201, width = 304
			name: "test8",
			item: rect.Item{Format: rect.Format{Width: 201, Height: 100}, Distance: rect.Format{Width: 3, Height: 1}},
			out:  rect.Format{Width: 304, Height: 201},
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 304, Height: 201},
				Fragments: base.Fragments{{ByWidth: 1, ByHeight: 2}, {ByWidth: 1, ByHeight: 1}},
				Total:     3,
				RestArea:  0,
			},
		},
		{
			name: "test9",
			item: rect.Item{Format: rect.Format{Width: 201, Height: 100}, Distance: rect.Format{Width: 1, Height: 1}},
			out:  rect.Format{Width: 302, Height: 201},
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 302, Height: 201},
				Fragments: base.Fragments{{ByWidth: 1, ByHeight: 2}, {ByWidth: 1, ByHeight: 1}},
				Total:     3,
				RestArea:  0,
			},
		},
		{
			// здесь расстояние между элементами не должно влиять
			name: "test10",
			item: rect.Item{Format: rect.Format{Width: 301, Height: 201}, Distance: rect.Format{Width: 10, Height: 10}},
			out:  rect.Format{Width: 302, Height: 201},
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 301, Height: 201},
				Fragments: base.Fragments{{ByWidth: 1, ByHeight: 1}},
				Total:     1,
				RestArea:  201,
			},
		},
		{
			name: "test11",
			item: rect.Item{Format: rect.Format{Width: 301, Height: 202}, Distance: rect.Format{Width: 10, Height: 10}},
			out:  rect.Format{Width: 302, Height: 201},
			want: imposition.AlgoResult{
				Layout:    rect.Format{},
				Fragments: nil,
				Total:     0,
				RestArea:  0,
			},
		},
		{
			// bug fixed
			name: "test12",
			item: rect.Item{Format: rect.Format{Width: 100, Height: 200}, Distance: rect.Format{Width: 2, Height: 2}},
			out:  rect.Format{Width: 442, Height: 304}, // 450x320 - frame 8x16
			want: imposition.AlgoResult{
				Layout:    rect.Format{Width: 406, Height: 302},
				Fragments: base.Fragments{{ByWidth: 4, ByHeight: 1}, {ByWidth: 1, ByHeight: 2}},
				Total:     6,
				RestArea:  12364,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			im := imposition.New(newLoggerStub(ctrl))

			got, err := im.Calc(tt.item, tt.out, imposition.Options{AllowRotation: true})
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}
