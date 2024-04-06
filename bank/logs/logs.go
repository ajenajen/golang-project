package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() { //init ทำงานก่อน main
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "" //ปิด stacktrace

	var err error
	log, err = config.Build(zap.AddCallerSkip(1)) // zap.AddCallerSkip(1) เพราะตอนนี้ caller อ่านว่า log จากไฟล์นี้ เลยต้อง skip
	if err != nil {
		panic(err)
	}

	// Log, _ = zap.NewProduction() // set config ของ zap เปลี่ยนไปมาได้
	// NewProduction จะ return เป็น json >> {"level":"info","ts":1712410785.4858408,"caller":"bank/main.go:66","msg":"Banking service started at port 8000"}
	// NewDevelopment จะ return เป็น console >>     INFO    bank/main.go:66 Banking service started at port 8000
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

func Error(message interface{}, fields ...zap.Field) { //interface คือรับทุก type
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), fields...)
	case string:
		log.Error(v, fields...)
	}
}
