package slices

import (
	"reflect"
	"sort"
	"strings"
)

// Difference returns the elements in `a` that aren't in `b`.
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// CompareMaps sorts and compares two slices containing map[string]interface{}
func CompareMaps(a, b []map[string]interface{}, sortKey string) bool {
	sort.Slice(a, func(i, j int) bool {
		id1 := a[i][sortKey].(string)
		id2 := a[j][sortKey].(string)
		return strings.ToUpper(id1) < strings.ToUpper(id2)
	})
	sort.Slice(b, func(i, j int) bool {
		id1 := b[i][sortKey].(string)
		id2 := b[j][sortKey].(string)
		return strings.ToUpper(id1) < strings.ToUpper(id2)
	})
	return reflect.DeepEqual(a, b)
}