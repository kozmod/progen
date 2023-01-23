package flag

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

//goland:noinspection SpellCheckingInspection
const (
	defaultConfigFilePath = "progen.yml"
)

func ParseFlags() entity.Flags {
	var f entity.Flags
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

	for _, arg := range flag.Args() {
		if strings.TrimSpace(arg) == entity.Dash {
			f.ReadStdin = true
			break
		}
	}

	return f
}
