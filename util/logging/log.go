package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"goprox/util/files"
	"log"
	"os"
	"path/filepath"
)

func LoggerSetup() *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(cfg)

	err := files.CreateDirectoryIfNotExists("logs", 0644)
	logFile, err := os.OpenFile(filepath.Join("logs", "goprox.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error when configuring zap: ", err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLevel := zapcore.DebugLevel
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, defaultLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger
}
