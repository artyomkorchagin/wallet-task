package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		validate    func(*zap.Logger, *testing.T)
		expectError bool
	}{
		{
			name: "production logger created successfully",
			validate: func(logger *zap.Logger, t *testing.T) {
				assert.NotNil(t, logger)

				buf := &bytes.Buffer{}
				encoderCfg := zap.NewProductionEncoderConfig()
				encoderCfg.TimeKey = "timestamp"
				encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

				core := zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderCfg),
					zapcore.AddSync(buf),
					zap.InfoLevel,
				)
				testLogger := zap.New(core)

				testLogger.Info("test message", zap.String("key", "value"))

				var log map[string]interface{}
				err := json.Unmarshal(buf.Bytes(), &log)
				assert.NoError(t, err)

				assert.Contains(t, log, "timestamp")
				assert.NotEmpty(t, log["timestamp"])

				assert.Equal(t, "info", log["level"])

				assert.Equal(t, "test message", log["msg"])
				assert.Equal(t, "value", log["key"])
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger()
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, logger)
			defer logger.Sync()

			if tt.validate != nil {
				tt.validate(logger, t)
			}
		})
	}
}

func TestNewDevelopmentLogger(t *testing.T) {
	tests := []struct {
		name        string
		validate    func(*zap.Logger, *testing.T)
		expectError bool
	}{
		{
			name: "development logger created successfully",
			validate: func(logger *zap.Logger, t *testing.T) {
				assert.NotNil(t, logger)

				observedZapCore, _ := observer.New(zap.DebugLevel)
				_ = zap.New(observedZapCore)

				buf := &bytes.Buffer{}
				encoderCfg := zap.NewDevelopmentEncoderConfig()
				encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

				core := zapcore.NewCore(
					zapcore.NewConsoleEncoder(encoderCfg),
					zapcore.AddSync(buf),
					zap.DebugLevel,
				)
				testLogger := zap.New(core)

				testLogger.Info("test dev message")

				output := buf.String()
				assert.Contains(t, output, "INFO")
				assert.Contains(t, output, "test dev message")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewDevelopmentLogger()
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, logger)
			defer logger.Sync()

			if tt.validate != nil {
				tt.validate(logger, t)
			}
		})
	}
}

func TestLoggers_AreUsable(t *testing.T) {
	prodLogger, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, prodLogger)
	prodLogger.Info("test production log")
	prodLogger.Sync()

	devLogger, err := NewDevelopmentLogger()
	assert.NoError(t, err)
	assert.NotNil(t, devLogger)
	devLogger.Info("test development log")
	devLogger.Sync()
}
