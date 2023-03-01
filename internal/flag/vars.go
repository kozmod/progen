package flag

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

var (
	ErrVariableNotSet = fmt.Errorf("value not set: mast be separated by [%s]", entity.EqualsSign)
)

type TemplateVarsFlag struct {
	Vars map[string]any
}

func (v *TemplateVarsFlag) String() string {
	return fmt.Sprintf("%v", v.Vars)
}

func (v *TemplateVarsFlag) Set(s string) error {
	if v.Vars == nil {
		v.Vars = make(map[string]any)
	}

	splitVar := strings.Split(s, entity.Dot)
	current := v.Vars
	if len(splitVar) > 0 {
		for i, key := range splitVar {
			key = strings.TrimSpace(key)
			if key == entity.Empty {
				continue
			}

			if i != len(splitVar)-1 {
				val, exists := v.Vars[key]
				if !exists {
					next := make(map[string]any, 1)
					current[key] = next
					current = next
					continue
				}
				if m, ok := val.(map[string]any); ok {
					current = m
					continue
				}
			}

			keyVal := strings.SplitN(key, entity.EqualsSign, 2)
			if len(keyVal) < 2 {
				return xerrors.Errorf("%w", ErrVariableNotSet)
			}

			current[keyVal[0]] = keyVal[1]
		}
	}

	return nil
}
