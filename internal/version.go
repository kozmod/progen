package internal

import (
	"runtime/debug"
	"strings"
	"time"

	"github.com/kozmod/progen/internal/entity"
)

const (
	_defaultVersion    = "v0.0.0"
	_sourceBuildSuffix = ".src"

	_debugRevisionKey = "vcs.revision"
	_debugTimeKey     = "vcs.time"
	_debugTimeLayout  = "20060102150405"
	_develVersion     = "(devel)"

	_versionSeparator = entity.Dash
	_empty            = entity.Empty
)

var (
	Version string
)

func GetVersion() string {
	if v := strings.TrimSpace(Version); Version != entity.Empty {
		return v
	}

	version, ok := debug.ReadBuildInfo()
	if ok && version.Main.Version != _develVersion && version.Main.Version != _empty {
		return version.Main.Version
	}

	var sb strings.Builder
	sb.WriteString(_defaultVersion)

	if ok {
		values := make(map[int]string, 2)
		for _, setting := range version.Settings {
			key, value := setting.Key, strings.TrimSpace(setting.Value)
			if value == entity.Empty {
				continue
			}

			switch key {
			case _debugTimeKey:
				versionTime, err := time.Parse(time.RFC3339, value)
				if err != nil {
					continue
				}
				values[1] = versionTime.Format(_debugTimeLayout)
			case _debugRevisionKey:
				values[2] = value
			}
		}

		for i := 1; i <= 2; i++ {
			value, ok := values[i]
			if !ok {
				continue
			}
			sb.WriteString(_versionSeparator)
			sb.WriteString(value)
		}
	}

	sb.WriteString(_sourceBuildSuffix)
	return sb.String()
}
