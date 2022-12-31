package config

import "github.com/kozmod/progen/interanl/proc"

func MustConfigureMkdirProc(conf Config) (proc.Proc, error) {
	if len(conf.Dirs) == 0 {
		return nil, nil
	}
	return proc.NewMkdirAllProc(uniqueVal(conf.Dirs)), nil
}

func uniqueVal(in []string) []string {
	set := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
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
