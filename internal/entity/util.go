package entity

import "reflect"

func MergeKeys(dst, src map[string]any) map[string]any {
	if dst == nil {
		return src
	}
	for key, rightVal := range src {
		if leftVal, present := dst[key]; present {
			switch rightVal.(type) {
			case map[string]any:
				dst[key] = MergeKeys(leftVal.(map[string]any), rightVal.(map[string]any))
			default:
				dst[key] = rightVal
			}
		} else {
			dst[key] = rightVal
		}
	}
	return dst
}

func NotNilValues(values ...any) int {
	counter := 0
	for _, value := range values {
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}
		if value == nil {
			continue
		}
		counter++
	}
	return counter
}
