package zaplog

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// set Logger as global variable
var Logger *zap.Logger

func init() {
	encoderJSON := encoder("JSON")
	encoderDefault := encoder("DEFAULT")

	core := zapcore.NewTee(
		zapcore.NewCore(encoderDefault, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(encoderJSON, logWriter("GIN-PRODUCTION.log"), zap.InfoLevel),
		zapcore.NewCore(encoderJSON, logWriter("GIN-PRODUCTION-ERROR.log"), zap.ErrorLevel),
	)

	Logger = zap.New(core)
	defer Logger.Sync()
}

func encoder(output string) zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	// https://github.com/uber-go/zap/blob/master/zapcore/encoder.go
	// config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05")

	switch output {
	case "JSON":
		return zapcore.NewJSONEncoder(config)
	case "DEFAULT":
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewConsoleEncoder(config)
	}
}

func logWriter(fname string) zapcore.WriteSyncer {
	lumberjack := &lumberjack.Logger{
		Filename:   fname,
		MaxSize:    1000, // megabytes
		MaxBackups: 10,
		MaxAge:     90, // days
	}
	return zapcore.AddSync(lumberjack)
}

func TimeTrack(start time.Time, fn string, name string) {
	elapsed := time.Since(start)
	Logger.Info(
		fmt.Sprintf("[GIN|%s] %v | %s | %s", fn, time.Now().Format("2006/01/02 - 15:04:05"),
			name,
			elapsed,
		),
	)
}
