package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"print-shop-back/internal/warehousing/util"
)

func TestPrepareTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags []string
		want []string
	}{
		{
			name: "test1",
			tags: nil,
			want: nil,
		},
		{
			name: "test2",
			tags: []string{},
			want: nil,
		},
		{
			name: "test3",
			tags: []string{"tag1", "tag1"},
			want: []string{"TAG1"},
		},
		{
			name: "test4",
			tags: []string{"tag2", "tag1"},
			want: []string{"TAG1", "TAG2"},
		},
		{
			name: "test5",
			tags: []string{" ", " \n ", "\n\t"},
			want: nil,
		},
		{
			name: "test6",
			tags: []string{" ", " tag4", " \n ", "\n\t", "TAG1 ", " tag5 ", "   ", " tag2 ", "tag4 "},
			want: []string{"TAG1", "TAG2", "TAG4", "TAG5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.PrepareTags(tt.tags)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrepareTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "test1",
			tag:  "",
			want: "",
		},
		{
			name: "test3",
			tag:  "tag1",
			want: "TAG1",
		},
		{
			name: "test4",
			tag:  "  tag1  ",
			want: "TAG1",
		},
		{
			name: "test7",
			tag:  "\n\t",
			want: "",
		},
		{
			name: "test8",
			tag:  " \n \t",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.PrepareTag(tt.tag)

			assert.Equal(t, tt.want, got)
		})
	}
}
