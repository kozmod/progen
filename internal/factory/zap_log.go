package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"
)

func NewLogger(verbose bool) (entity.LoggerWrapper, error) {
	lvl := zap.ErrorLevel
	if verbose {
		lvl = zap.InfoLevel
	}

	atomicLvl := zap.NewAtomicLevelAt(lvl)

	cfg := zap.Config{
		Level:       atomicLvl,
		Development: false,
		Sampling:    nil,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   nil,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	base, err := cfg.Build()
	if err != nil {
		return nil, xerrors.Errorf("create new logger: %w", err)
	}

	return &zapLoggerWrapper{
		SugaredLogger: base.Sugar(),
		AtomicLevel:   &atomicLvl,
		initLvl:       lvl,
	}, nil
}

type zapLoggerWrapper struct {
	*zap.SugaredLogger
	*zap.AtomicLevel
	initLvl zapcore.Level
}

func (lw *zapLoggerWrapper) ForceInfof(template string, args ...interface{}) {
	lw.TrySetInfoLevel()
	lw.Infof(template, args...)
	lw.TrySetInitLevel()
}

func (lw *zapLoggerWrapper) TrySetInfoLevel() {
	if lw.AtomicLevel.Level() != zap.InfoLevel {
		lw.SetLevel(zap.InfoLevel)
	}
}

func (lw *zapLoggerWrapper) TrySetInitLevel() {
	if lw.AtomicLevel.Level() != lw.initLvl {
		lw.SetLevel(lw.initLvl)
	}
}
