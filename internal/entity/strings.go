package entity

import "strings"

type StringsFn struct{}

func (sfn StringsFn) Replace(input, from, to string, n int) string {
	return strings.Replace(input, from, to, n)
}
