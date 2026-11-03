package utils

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"justcallmesu.com/golang-base/internal/common/apperror"
)

func ResolveSort(allowed map[string]string, key string, dir string) (string, string, error) {
	var matchedColumn string
	var ok bool
	var dirUpper = strings.ToUpper(dir)

	if matchedColumn, ok = allowed[key]; !ok {
		return "", "", apperror.BadRequest(fmt.Sprintf("invalid sort key, allowed keys are: %v", strings.Join(slices.Sorted(maps.Keys(allowed)), ", ")))
	}

	if dirUpper != "ASC" && dirUpper != "DESC" {
		return "", "", apperror.BadRequest("invalid sort direction, allowed directions are: ASC, DESC")
	}

	return matchedColumn, dirUpper, nil
}
