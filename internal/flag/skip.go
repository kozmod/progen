package flag

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type SkipFlag []string

func (s *SkipFlag) String() string {
	type alias struct {
		skip []string
	}
	a := alias{skip: *s}
	return fmt.Sprintf("%v", a.skip)
}

func (s *SkipFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == entity.Empty {
		return nil
	}
	*s = append(*s, value)
	return nil
}
