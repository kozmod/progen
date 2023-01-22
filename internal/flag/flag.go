package flag

import (
	"flag"
	"fmt"
)

//goland:noinspection SpellCheckingInspection
const (
	defaultConfigFilePath = "progen.yml"
)

type Flags struct {
	ConfigPath string
	Verbose    bool
	DryRun     bool
	Version    bool
}

func ParseFlags() Flags {
	var f Flags
	flag.StringVar(
		&f.ConfigPath,
		"f",
		defaultConfigFilePath,
		fmt.Sprintf("configuration file path (default file: %s)", defaultConfigFilePath))
	flag.BoolVar(
		&f.Verbose,
		"v",
		false,
		"verbose output")
	flag.BoolVar(
		&f.DryRun,
		"dr",
		false,
		"dry run mode (can be combine with `-v`)")
	flag.BoolVar(
		&f.Version,
		"version",
		false,
		"output version")
	flag.Parse()
	return f
}
