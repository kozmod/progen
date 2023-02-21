package entity

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unique(t *testing.T) {
	t.Parallel()

	type useCase[T comparable] struct {
		in  []T
		exp []T
	}

	stringCases := []useCase[string]{
		{
			in:  []string{"a", "a", "a"},
			exp: []string{"a"},
		},
		{
			in:  []string{"a", "b", "c"},
			exp: []string{"a", "b", "c"},
		},
		{
			in:  []string{},
			exp: []string{},
		},
		{
			in:  nil,
			exp: []string{},
		},
	}

	for i, tc := range stringCases {
		t.Run(fmt.Sprintf("str_case_%d", i), func(t *testing.T) {
			res := Unique(tc.in)
			assert.Equal(t, tc.exp, res)
		})
	}

	intCases := []useCase[int32]{
		{
			in:  []int32{1, 1, 1},
			exp: []int32{1},
		},
		{
			in:  []int32{1, 2, 3},
			exp: []int32{1, 2, 3},
		},
		{
			in:  []int32{},
			exp: []int32{},
		},
		{
			in:  nil,
			exp: []int32{},
		},
	}

	for i, tc := range intCases {
		t.Run(fmt.Sprintf("int_case_%d", i), func(t *testing.T) {
			res := Unique(tc.in)
			assert.Equal(t, tc.exp, res)
		})
	}
}

func Test_MergeKeys(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		src,
		dst,
		exp map[string]any
	}{
		{
			src: map[string]any{},
			dst: map[string]any{"var": map[string]any{"Y": "SOME_DST"}},
			exp: map[string]any{"var": map[string]any{"Y": "SOME_DST"}},
		},
		{
			src: nil,
			dst: map[string]any{"var": map[string]any{"Y": "SOME_DST"}},
			exp: map[string]any{"var": map[string]any{"Y": "SOME_DST"}},
		},
		{
			src: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
			dst: map[string]any{},
			exp: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
		},
		{
			src: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
			dst: nil,
			exp: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
		},
		{
			src: map[string]any{},
			dst: map[string]any{},
			exp: map[string]any{},
		},
		{
			src: nil,
			dst: nil,
			exp: nil,
		},
		{
			src: map[string]any{"var": "FROM_SRC"},
			dst: map[string]any{"var": "SOME_DST_2"},
			exp: map[string]any{"var": "FROM_SRC"},
		},
		{
			src: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
			dst: map[string]any{"var": map[string]any{"Y": "SOME_DST", "X": "SOME_DST_2"}},
			exp: map[string]any{"var": map[string]any{"X": "FROM_SRC", "Y": "SOME_DST"}},
		},
		{
			src: map[string]any{"var": map[string]any{"X": "FROM_SRC"}},
			dst: map[string]any{"var": map[string]any{"Y": "SOME_DST"}},
			exp: map[string]any{"var": map[string]any{"X": "FROM_SRC", "Y": "SOME_DST"}},
		},
		{
			src: map[string]any{"var": "FROM_SRC"},
			dst: map[string]any{"var2": map[string]any{"Y": "SOME_DST"}},
			exp: map[string]any{
				"var":  "FROM_SRC",
				"var2": map[string]any{"Y": "SOME_DST"},
			},
		},
	}

	for i, testCase := range testCases {
		tc := testCase
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			res := MergeKeys(tc.dst, tc.src)
			assert.Truef(t, reflect.DeepEqual(tc.exp, res), fmt.Sprintf("expected: %v, result: %v", tc.exp, res))
		})
	}
}

func Test_NotNilValues(t *testing.T) {
	t.Parallel()

	var (
		nilStr     *string
		defaultStr string

		nilInt     *int
		defaultInt int

		nilStrict     *struct{}
		defaultStruct struct{}

		ptrAny   *any
		valueAny any
	)

	test := []struct {
		in  []any
		exp int
	}{
		{in: []any{nilStr, nilInt, nilStrict, ptrAny, valueAny}, exp: 0},
		{in: []any{ptrAny}, exp: 0},
		{in: []any{valueAny}, exp: 0},
		{in: []any{defaultStr, defaultInt, defaultStruct}, exp: 3},
		{in: []any{nilStr, nilInt, nilStrict, defaultStr, defaultInt, defaultStruct}, exp: 3},
		{in: []any{nilInt, defaultInt}, exp: 1},
		{in: []any{}, exp: 0},
	}

	for i, tc := range test {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			res := NotNilValues(tc.in...)
			assert.Equalf(t, tc.exp, res, "case_%d", i)
		})
	}
}
