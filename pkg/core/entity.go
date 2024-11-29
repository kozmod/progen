package core

import (
	"github.com/kozmod/progen/internal/entity"
)

type (
	Action interface {
		add(engin *Engin)
	}

	File struct {
		Path string
		Data []byte
	}

	Cmd entity.Command
)

type Files entity.Action[[]File]

func (c Files) add(e *Engin) {
	if e != nil {
		files := convert(c.Val, func(s File) entity.UndefinedFile {
			return entity.UndefinedFile{
				Path: s.Path,
				Data: &s.Data,
			}
		})
		e.files = append(e.files, entity.Action[[]entity.UndefinedFile]{
			Name:     c.Name,
			Val:      files,
			Priority: c.Priority,
		})
	}
}

func FilesAction(name string, files ...File) Files {
	return Files{
		Name: name,
		Val:  files,
	}
}

type Command entity.Action[[]Cmd]

func (c Command) add(e *Engin) {
	if e != nil {
		commands := convert(c.Val, func(s Cmd) entity.Command {
			return entity.Command{
				Args: s.Args,
				Dir:  s.Dir,
				Cmd:  s.Cmd,
			}
		})
		e.cmd = append(e.cmd, entity.Action[[]entity.Command]{
			Name:     c.Name,
			Val:      commands,
			Priority: c.Priority,
		})
	}
}

func (c Command) WithPriority(priority int) Command {
	c.Priority = priority
	return c
}

func CmdAction(name string, commands ...Cmd) Command {
	return Command{
		Name: name,
		Val:  commands,
	}
}

type Dirs entity.Action[[]string]

func (d Dirs) add(e *Engin) {
	if e != nil {
		e.dirs = append(e.dirs, entity.Action[[]string]{
			Name:     d.Name,
			Val:      d.Val,
			Priority: d.Priority,
		})
	}
}

func (d Dirs) WithPriority(priority int) Dirs {
	d.Priority = priority
	return d
}

func DirsAction(name string, dirs ...string) Dirs {
	return Dirs{
		Name: name,
		Val:  dirs,
	}
}

type Rm entity.Action[[]string]

func (r Rm) add(e *Engin) {
	if e != nil {
		e.rm = append(e.rm, entity.Action[[]string]{
			Name:     r.Name,
			Val:      r.Val,
			Priority: r.Priority,
		})
	}
}

func (r Rm) WithPriority(priority int) Rm {
	r.Priority = priority
	return r
}

func RmAction(name string, rm ...string) Rm {
	return Rm{
		Name: name,
		Val:  rm,
	}
}

type Fs entity.Action[[]string]

func (f Fs) add(e *Engin) {
	if e != nil {
		e.fs = append(e.fs, entity.Action[[]string]{
			Name:     f.Name,
			Val:      f.Val,
			Priority: f.Priority,
		})
	}
}

func (f Fs) WithPriority(priority int) Fs {
	f.Priority = priority
	return f
}

func FsAction(name string, fs ...string) Fs {
	return Fs{
		Name: name,
		Val:  fs,
	}
}

func convert[S any, T any](s []S, fn func(s S) T) []T {
	res := make([]T, len(s))
	for i, val := range s {
		res[i] = fn(val)
	}
	return res
}
