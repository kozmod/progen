package core

import (
	"flag"
	"fmt"
	"os"

	internalFlags "github.com/kozmod/progen/internal/flag"
)

// Config for the [Engin].
type Config internalFlags.DefaultFlags

// AddConfigFlagSet adds flags to [*flag.FlagSet] and creates [Engin] config.
func AddConfigFlagSet(fs *flag.FlagSet) *Config {
	f := internalFlags.NewDefaultFlags(fs)
	return (*Config)(f)
}

// ParseFlags - parse flags and creates [Engin] config from [*flag.FlagSet].
func ParseFlags() (*Config, error) {
	var (
		args    = os.Args
		flagSet = flag.NewFlagSet(args[0], flag.ExitOnError)
		f       = AddConfigFlagSet(flagSet)
		err     = flagSet.Parse(args[1:])
	)
	if err != nil {
		return f, fmt.Errorf("parse config from flags: %w", err)
	}
	return f, nil
}
