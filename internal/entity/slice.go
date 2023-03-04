package entity

type SLiceFn struct{}

func (sfn SLiceFn) New(elems ...any) []any {
	return elems
}

func (sfn SLiceFn) Append(s []any, elems ...any) []any {
	return append(s, elems...)
}
