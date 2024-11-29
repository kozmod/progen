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
	flagKeyPrintConfig                 = "printconf"
	flagKeyDryRun                      = "dr"
	flagKeyVersion                     = "version"
	flagKeyTemplateVariables           = "tvar"
	flagKeyApplicationWorkingDirectory = "awd"
	flagKeySkip                        = "skip"
	flagKeyPreprocessingAllFiles       = "pf"
	flagKeyMissingKey                  = "missingkey"
	flagKeyGroup                       = "gp"
)

var (
	ErrDashFlagNotLast = fmt.Errorf(`'-' must be the last argument`)
)

type DefaultFlags struct {
	Verbose              bool
	DryRun               bool
	TemplateVars         TemplateVarsFlag
	MissingKey           MissingKeyFlag
	PrintErrorStackTrace bool
}

type Flags struct {
	DefaultFlags

	ConfigPath           string
	Version              bool
	ReadStdin            bool
	AWD                  string // AWD application working directory
	Skip                 SkipFlag
	PreprocessFiles      bool
	Group                GroupFlag
	PrintProcessedConfig bool
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
	var (
		dFlags DefaultFlags
	)

	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	err := parseDefaultFlags(&dFlags, flagSet)
	if err != nil {
		log.Fatalf("parse default flags: %v", err)
	}

	var (
		flags = Flags{
			DefaultFlags: dFlags,
		}
	)

	err = parseAdditionalFlags(&flags, flagSet, args[1:])
	if err != nil {
		log.Fatalf("parse additioanl flags: %v", err)
	}
	return flags
}

func ParseDefault() DefaultFlags {
	args := os.Args
	var (
		flags DefaultFlags
	)
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)
	err := parseDefaultFlags(&flags, flagSet)
	if err != nil {
		log.Fatalf("parse default flags: %v", err)
	}
	return flags
}

func parseDefaultFlags(f *DefaultFlags, fs *flag.FlagSet) error {
	fs.BoolVar(
		&f.Verbose,
		flagKeyVarbose,
		false,
		"verbose output")
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
	fs.Var(
		&f.TemplateVars,
		flagKeyTemplateVariables,
		"template variables")
	return nil
}

func parseAdditionalFlags(f *Flags, fs *flag.FlagSet, args []string) error {
	fs.StringVar(
		&f.ConfigPath,
		flagKeyConfigFile,
		defaultConfigFilePath,
		"configuration file path")
	fs.BoolVar(
		&f.PrintProcessedConfig,
		flagKeyPrintConfig,
		false,
		"output processed config")
	fs.BoolVar(
		&f.Version,
		flagKeyVersion,
		false,
		"output version")
	fs.StringVar(
		&f.AWD,
		flagKeyApplicationWorkingDirectory,
		entity.Dot,
		"application working directory")
	fs.Var(
		&f.Skip,
		flagKeySkip,
		"list of skipping actions")
	fs.BoolVar(
		&f.PreprocessFiles,
		flagKeyPreprocessingAllFiles,
		true,
		"preprocessing all files before saving")
	fs.Var(
		&f.Group,
		flagKeyGroup,
		"list of executing groups",
	)
	err := fs.Parse(args)
	if err != nil {
		return xerrors.Errorf("parse args: %w", err)
	}

	for i, arg := range args {
		if strings.TrimSpace(arg) == entity.Dash {
			if i == len(args)-1 || (i < len(args)-1 && args[i+1] == entity.LessThan) {
				f.ReadStdin = true
				break
			}
			return xerrors.Errorf("%w", ErrDashFlagNotLast)
		}
	}

	return nil
}
