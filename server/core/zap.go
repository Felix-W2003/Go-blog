package core

import (
	"log"
	"os"
	"server/global"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// getLogWriter 返回一个 zapcore.WriteSyncer，该写入器利用 lumberjack 包，实现日志的滚动记录
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件的位置
		MaxSize:    maxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		MaxAge:     maxAge,     // 保留旧文件的最大天数
	}
	return zapcore.AddSync(lumberJackLogger)
}

// InitLogger 初始化并返回一个基于配置设置的新 zap.Logger 实例
func InitLogger() *zap.Logger {
	zapCfg := global.Config.Zap

	// 创建一个用于日志输出的 writeSyncer
	writeSyncer := getLogWriter(zapCfg.Filename, zapCfg.MaxSize, zapCfg.MaxBackups, zapCfg.MaxAge)

	// 如果配置了控制台输出，则添加控制台输出
	if zapCfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	// 创建日志格式化的编码器
	encoder := getEncoder()

	// 根据配置确定日志级别
	var logLevel zapcore.Level

	if err := logLevel.UnmarshalText([]byte(zapCfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	// 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())

	return logger
}

// getEncoder 返回一个为生产日志配置的 JSON 编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

/*整套的逻辑：读取yaml配置，组装zap+lumberjack，生成全局日志对象
zap负责日志格式、级别；lumberjack负责日志文件切割轮转

整体流程

InitLogger()入口函数执行顺序：
读取全局配置global.Config.Zap（从 config.yaml 解析出来的 zap 配置）
调用getLogWriter()创建文件输出（带日志切割）
如果开启控制台打印，就同时输出到文件 + 终端 stdout
调用getEncoder()创建 JSON 格式编码器（日志输出格式）
解析配置文件中的日志级别字符串 (info/debug/warn) 转为 zap 内部 Level 类型
组装zapcore.Core，构建*zap.Logger，开启行号调用者追踪 zap.AddCaller()
返回 logger 实例给全局变量保存，整个项目就用这个 logger 打日志。
*/
