package flag

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type GroupFlag []string

func (s *GroupFlag) String() string {
	type alias struct {
		skip []string
	}
	a := alias{skip: *s}
	return fmt.Sprintf("%v", a.skip)
}

func (s *GroupFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == entity.Empty {
		return nil
	}
	*s = append(*s, value)
	return nil
}
