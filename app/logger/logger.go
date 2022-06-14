package logger

import (
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - holds a logger object
type Logger struct {
	Loggers map[string]*zap.Logger
}

var singleton *Logger
var once sync.Once

// Get returns the singleton logger instance
func Get(name string) *zap.Logger {
	once.Do(func() {
		singleton = &Logger{Loggers: make(map[string]*zap.Logger)}
	})

	logger, ok := singleton.Loggers[name]
	if !ok {
		// _ = os.Mkdir("logs", 0700)
		// logpath := path.Join("logs", name+".log")
		var err error
		config := zap.Config{
			Encoding:    "console",
			Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
			OutputPaths: []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:     "time",
				EncodeTime:  zapcore.RFC3339TimeEncoder,
				LevelKey:    "level",
				EncodeLevel: zapcore.CapitalColorLevelEncoder,
				MessageKey:  "message",
				NameKey:     "name",
				EncodeName:  zapcore.FullNameEncoder,
			},
		}

		if os.Getenv("ENV") == "prod" {
			// Only log warnings
			config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		}

		logger, err = config.Build()

		if err != nil {
			log.Fatalln("Error loading logger:", err)
		}

		singleton.Loggers[name] = logger.Named(name)
		logger = singleton.Loggers[name]
	}
	return logger
}
