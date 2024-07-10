package entity

type AppendVPlusOrV func(s string) string

func NewAppendVPlusOrV(vPlus bool) AppendVPlusOrV {
	return func(s string) string {
		const vp, v = "%+v", "%v"
		if vPlus {
			return s + vp
		}
		return s + v
	}
}
