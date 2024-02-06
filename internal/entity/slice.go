package entity

type SliceFn struct{}

func (sfn SliceFn) New(elems ...any) []any {
	return elems
}

func (sfn SliceFn) Append(s []any, elems ...any) []any {
	return append(s, elems...)
}
