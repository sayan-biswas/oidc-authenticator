package logger

import (
	"os"
	"strconv"

	"github.com/mattn/go-colorable"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Level zap.AtomicLevel

func init() {

	var encoder zapcore.Encoder
	var encoderConfig zapcore.EncoderConfig
	var option []zap.Option
	var core zapcore.Core

	debug, _ := strconv.ParseBool(pflag.Lookup("debug").Value.String())
	logLevel := pflag.Lookup("logLevel").Value.String()
	logOutput := pflag.Lookup("logOutput").Value.String()

	DEBUG, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if debug || DEBUG {
		logLevel = "debug"
		pflag.Set("debug", "true")
		pflag.Set("logLevel", "debug")
	}

	Level = zap.NewAtomicLevel()

	encoderConfig = zapcore.EncoderConfig{
		LevelKey:       "level",
		MessageKey:     "message",
		TimeKey:        "time",
		NameKey:        zapcore.OmitKey,
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	switch logOutput {
	case "json":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core = zapcore.NewCore(
		encoder,
		zapcore.AddSync(colorable.NewColorableStdout()),
		Level,
	)

	option = []zap.Option{
		zap.ErrorOutput(os.Stderr),
		zap.AddCaller(),
	}

	Logger = zap.New(core, option...)
	SetLevel(logLevel)
}

func SetLevel(level string) error {
	var logLevel zapcore.Level
	Logger.Info("Changing log level", zap.String("level", level))
	if setError := logLevel.Set(level); setError != nil {
		Logger.Error("Error changing log level", zap.Error(setError))
		return setError
	}
	Level.SetLevel(logLevel)
	return nil
}

func Get() *zap.Logger {
	return Logger
}
