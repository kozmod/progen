package internal

import "runtime/debug"

const (
	_defaultVersion = "v0.0.0-local"
)

var Version string

func GetVersion() string {
	if Version != "" {
		return Version
	}
	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != "(devel)" && version.Main.Version != "" {
		return version.Main.Version
	}
	return _defaultVersion
}
