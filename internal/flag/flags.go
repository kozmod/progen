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

type Flags struct {
	ConfigPath   string
	Verbose      bool
	DryRun       bool
	Version      bool
	ReadStdin    bool
	TemplateVars TemplateVarsFlag
	// AWD application working directory
	AWD  string
	Skip SkipFlag
}

func (f *Flags) FileLocationMessage() string {
	switch {
	case f == nil:
		return entity.Empty
	case f.ReadStdin:
		return "stdin"
	case strings.TrimSpace(f.ConfigPath) != entity.Empty:
		return f.ConfigPath
	default:
		return entity.Empty
	}
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
	//goland:noinspection SpellCheckingInspection
	flag.Var(
		&f.TemplateVars,
		"tvar",
		"template variables (override config variables tree)")
	flag.StringVar(
		&f.AWD,
		"awd",
		entity.Dot,
		"application working directory")
	flag.Var(
		&f.Skip,
		"skip",
		"list of skipping yaml tags")
	flag.Parse()

	for _, arg := range flag.Args() {
		if strings.TrimSpace(arg) == entity.Dash {
			f.ReadStdin = true
			break
		}
	}

	return f
}
