package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unique(t *testing.T) {
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
