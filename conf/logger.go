package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() *zap.SugaredLogger {
	logMode := zapcore.DebugLevel
	if !viper.GetBool("mode.dev") {
		logMode = zapcore.InfoLevel
	}
	core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriterSyncer(), zapcore.AddSync(os.Stdout)), logMode)
	return zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncer() zapcore.WriteSyncer {
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly)
	fmt.Println(stLogFilePath)

	lumberjackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("log.MaxSize"),
		MaxAge:     viper.GetInt("log.MaxAge"),
		MaxBackups: viper.GetInt("log.MaxBackups"),
		LocalTime:  false,
		Compress:   false,
	}

	return zapcore.AddSync(lumberjackSyncer)
}
