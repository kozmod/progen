package entity

import (
	"hash/maphash"
	"math/rand"
	"unicode"
)

// Default set, matches "[a-zA-Z0-9_.-]"
const (
	_lettersAlpha    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	_lettersAlphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	_letterIdxBits = 6                     // 6 bits to represent a letter index
	_letterIdxMask = 1<<_letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	_letterIdxMax  = 63 / _letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	_lettersASCII string

	mapHashSrc = rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
)

func init() {
	var (
		lower = int32(Space[0])
		upper = int32(Tilda[0])
	)

	chars := make([]rune, 0, int(upper)-int(lower))
	for r := lower; r <= upper; r++ {
		if unicode.IsGraphic(r) {
			chars = append(chars, r)
		}
	}
	_lettersASCII = string(chars)
}

// RandomFn has to generate random string value
type RandomFn struct{}

// Alpha Generates a random alphabetical (A-Z, a-z) string of a desired length.
func (RandomFn) Alpha(n int) string {
	return randomString(n, _lettersAlpha)
}

// AlphaNum Generates a random alphanumeric (0-9, A-Z, a-z) string of a desired length.
func (RandomFn) AlphaNum(n int) string {
	return randomString(n, _lettersAlphaNum)
}

// ASCII Generates a random string of a desired length, containing the set of printable characters from the 7-bit ASCII set.
// This includes space (’ ‘), but no other whitespace character
func (RandomFn) ASCII(n int) string {
	return randomString(n, _lettersASCII)
}

func randomString(n int, set string) string {
	src := mapHashSrc
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), _letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), _letterIdxMax
		}
		if idx := int(cache & _letterIdxMask); idx < len(set) {
			b[i] = set[idx]
			i--
		}
		cache >>= _letterIdxBits
		remain--
	}

	return string(b)
}
