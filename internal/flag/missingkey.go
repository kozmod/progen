package flag

import (
	"fmt"
	"github.com/kozmod/progen/internal/entity"
	"strings"
)

type MissingKeyFlag string

func (s *MissingKeyFlag) String() string {
	type alias struct {
		mk string
	}

	var a alias
	switch {
	case s == nil,
		strings.TrimSpace(string(*s)) == entity.Empty:
		a = alias{mk: string(entity.MissingKeyError)}
	default:
		a = alias{mk: string(*s)}
	}
	return fmt.Sprintf("%v=%v", entity.TemplateOptionsMissingKey, a.mk)
}

func (s *MissingKeyFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == entity.Empty {
		*s = MissingKeyFlag(entity.MissingKeyError)
		return nil
	}

	if err := entity.MissingKeyValue(value).Valid(); err != nil {
		return err
	}

	*s = MissingKeyFlag(value)
	return nil
}
