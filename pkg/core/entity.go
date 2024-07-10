package core

import (
	"github.com/kozmod/progen/internal/entity"
)

type (
	File struct {
		Path string
		Data []byte
	}

	Cmd entity.Command
)

func toEntityActions[S ActionVal, T any](convert func(s S) T, in ...entity.Action[[]S]) []entity.Action[[]T] {
	res := make([]entity.Action[[]T], len(in))
	for i, action := range in {
		target := make([]T, len(in))
		for j, s := range action.Val {
			target[j] = convert(s)
		}
		res[i] = entity.Action[[]T]{
			Priority: action.Priority,
			Name:     action.Name,
			Val:      target,
		}
	}
	return res
}

type ActionVal interface {
	string | Cmd | File
}

func NewAction[T ActionVal](name string, val ...T) entity.Action[[]T] {
	return entity.Action[[]T]{
		Name:     name,
		Priority: 0,
		Val:      val,
	}
}

func WithPriority[T any](action entity.Action[[]T], priority int) entity.Action[[]T] {
	action.Priority = priority
	return action
}
