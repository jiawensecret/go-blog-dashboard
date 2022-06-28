package core

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"platform/global"
	"platform/utils"
	"time"
)

var level zapcore.Level

func Zap() *zap.SugaredLogger {
	if ok, _ := utils.PathExists(global.Config.System.Path); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(global.Config.System.Path, os.ModePerm)
	}

	switch global.Config.System.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	var encoderCfg zapcore.EncoderConfig
	if os.Getenv("SERVER_MODE") == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.TimeKey = "time"
	encoderCfg.NameKey = "name"
	encoderCfg.MessageKey = "message"
	encoderCfg.FunctionKey = "function"
	encoderCfg.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	if os.Getenv("LOG_ENCODING") == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	logWriter := getWriter(global.Config.System.Path)
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(level))

	var logger *zap.Logger
	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(level))
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	sugarLogger := logger.Sugar()
	if err := sugarLogger.Sync(); err != nil {
		sugarLogger.Error(err)
	}

	return sugarLogger
}

func getWriter(logPath string) zapcore.WriteSyncer {
	if logPath == "/dev/stdout" {
		return zapcore.Lock(os.Stdout)
	}
	if logPath == "/dev/stderr" {
		return zapcore.Lock(os.Stderr)
	}

	hook, err := rotatelogs.New(
		logPath+"-%Y-%m-%d"+".log",
		rotatelogs.WithLinkName(logPath+".log"),
		rotatelogs.WithMaxAge(time.Hour*24*30),    // 保存30天
		rotatelogs.WithRotationTime(time.Hour*24), //切割频率 24小时
	)
	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return zapcore.AddSync(hook)
}
