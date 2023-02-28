package flag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

//goland:noinspection SpellCheckingInspection
const (
	defaultConfigFilePath = "progen.yml"

	flagKeyConfigFile                  = "f"
	flagKeyVarbose                     = "v"
	flagKeyErrorStackTrace             = "errtrace"
	flagKeyDryRun                      = "dr"
	flagKeyVersion                     = "version"
	flagKeyTemplateVariables           = "tvar"
	flagKeyApplicationWorkingDirectory = "awd"
	flagKeySkip                        = "skip"
	flagKeyPreprocessingAllFiles       = "pf"
	flagKeyMissingKey                  = "missingkey"
)

var (
	ErrDashFlagNotLast = fmt.Errorf(`'-' must be the last argument`)
)

type Flags struct {
	ConfigPath           string
	Verbose              bool
	DryRun               bool
	Version              bool
	ReadStdin            bool
	TemplateVars         TemplateVarsFlag
	AWD                  string // AWD application working directory
	Skip                 SkipFlag
	PreprocessFiles      bool
	MissingKey           MissingKeyFlag
	PrintErrorStackTrace bool
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

func Parse() Flags {
	args := os.Args
	flags, err := parseFlags(flag.NewFlagSet(args[0], flag.ExitOnError), args[1:])
	if err != nil {
		log.Fatalf("parse flags: %v", err)
	}
	return flags
}

func parseFlags(fs *flag.FlagSet, args []string) (Flags, error) {
	var (
		f Flags
	)
	fs.StringVar(
		&f.ConfigPath,
		flagKeyConfigFile,
		defaultConfigFilePath,
		"configuration file path")
	fs.BoolVar(
		&f.Verbose,
		flagKeyVarbose,
		false,
		"verbose output")
	//goland:noinspection SpellCheckingInspection
	fs.BoolVar(
		&f.PrintErrorStackTrace,
		flagKeyErrorStackTrace,
		false,
		"output errors stacktrace")
	fs.BoolVar(
		&f.DryRun,
		flagKeyDryRun,
		false,
		"dry run mode (can be combine with `-v`)")
	fs.BoolVar(
		&f.Version,
		flagKeyVersion,
		false,
		"output version")
	fs.Var(
		&f.TemplateVars,
		flagKeyTemplateVariables,
		"template variables (override config variables tree)")
	fs.StringVar(
		&f.AWD,
		flagKeyApplicationWorkingDirectory,
		entity.Dot,
		"application working directory")
	fs.Var(
		&f.Skip,
		flagKeySkip,
		"list of skipping 'yaml' tags")
	fs.BoolVar(
		&f.PreprocessFiles,
		flagKeyPreprocessingAllFiles,
		true,
		"preprocessing all files before saving")
	fs.Var(
		&f.MissingKey,
		flagKeyMissingKey,
		fmt.Sprintf(
			"`missingkey` template option: %v, %v, %v, %v",
			entity.MissingKeyDefault,
			entity.MissingKeyInvalid,
			entity.MissingKeyZero,
			entity.MissingKeyError,
		),
	)
	err := fs.Parse(args)
	if err != nil {
		return f, err
	}

	for i, arg := range args {
		if strings.TrimSpace(arg) == entity.Dash {
			if i == len(args)-1 || (i < len(args)-1 && args[i+1] == entity.LessThan) {
				f.ReadStdin = true
				break
			}
			return f, xerrors.Errorf("%w", ErrDashFlagNotLast)
		}
	}

	return f, nil
}
