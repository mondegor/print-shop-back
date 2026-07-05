package util

import (
	"slices"
	"strings"
)

// PrepareTags - comment func.
func PrepareTags(tags []string) []string {
	j := 0

	// сдвиг непустых тегов вперёд, чтобы все пустые
	// теги остались в конце массива
	for i := 0; i < len(tags); i++ {
		if tags[j] = PrepareTag(tags[i]); tags[j] != "" {
			j++
		}
	}

	if j < len(tags) {
		clear(tags[j:])
		tags = tags[:j]
	}

	if len(tags) == 0 {
		return nil
	}

	slices.Sort(tags)

	return slices.Compact(tags)
}

// PrepareTag - comment func.
func PrepareTag(tag string) string {
	if tag = strings.TrimSpace(tag); tag == "" {
		return ""
	}

	return strings.ToUpper(tag)
}
