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
	*DefaultFlags

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
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)

	var (
		flags = Flags{
			DefaultFlags: NewDefaultFlags(flagSet),
		}
	)

	err := flags.Parse(flagSet, args[1:])
	if err != nil {
		log.Fatalf("parse additioanl flags: %v", err)
	}
	return flags
}

func NewDefaultFlags(fs *flag.FlagSet) *DefaultFlags {
	var f DefaultFlags
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
	return &f
}

func (f *DefaultFlags) Parse(fs *flag.FlagSet, args []string) error {
	err := fs.Parse(args)
	if err != nil {
		return xerrors.Errorf("parse args: %w", err)
	}
	return nil
}

func NewFlags(fs *flag.FlagSet) *Flags {
	var (
		f = Flags{DefaultFlags: NewDefaultFlags(fs)}
	)
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

	return &f
}

func (f *Flags) Parse(fs *flag.FlagSet, args []string) error {
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
