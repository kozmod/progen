package entity

func Unique[T comparable](in []T) []T {
	set := make(map[T]struct{}, len(in))
	out := make([]T, 0, len(in))
	for _, val := range in {
		_, ok := set[val]
		if ok {
			continue
		}
		out = append(out, val)
		set[val] = struct{}{}
	}
	return out
}
