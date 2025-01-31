package util

import (
	"golang.org/x/exp/constraints"
	"sort"
)

func SortResourcesById[T any](r []T, idGetter func(T) string) []T {
	sort.Slice(
		r, func(a, b int) bool {
			return idGetter(r[a]) < idGetter(r[b])
		},
	)
	return r
}

func GetSortedMapValues[T any](r map[string]T) []T {
	var keys []string
	for k := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result []T
	for _, k := range keys {
		result = append(result, r[k])
	}
	return result
}

func MaxBy[K constraints.Ordered, V any, P constraints.Ordered](r map[K]V, selector func(V) P) K {
	var currentMax K
	for k, v := range r {
		if _, ok := r[currentMax]; !ok {
			currentMax = k
		}
		if selector(r[currentMax]) < selector(v) || (selector(r[currentMax]) == selector(v) && currentMax < k) {
			currentMax = k
		}
	}
	return currentMax
}

func Contains[T comparable](list []T, a T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Keys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[K comparable, V any, V2 any](m map[K]V, mapper func(V) V2) map[K]V2 {
	result := map[K]V2{}
	for k, v := range m {
		result[k] = mapper(v)
	}
	return result
}

func MapSliceValues[V any, V2 any](values []V, mapper func(V) V2) []V2 {
	var result []V2
	for _, v := range values {
		result = append(result, mapper(v))
	}
	return result
}

func MapSliceValuesErr[V any, V2 any](values []V, mapper func(V) (V2, error)) ([]V2, error) {
	var result []V2
	for _, v := range values {
		mapped, err := mapper(v)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func FlatMapSliceValues[V any, V2 any](values []V, mapper func(V) []V2) []V2 {
	var result []V2
	for _, v := range values {
		for _, v2 := range mapper(v) {
			result = append(result, v2)
		}
	}
	return result
}
