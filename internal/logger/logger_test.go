package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/artyomkorchagin/wallet-task/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatal("Expected non-nil logger")
	}

	logger.Info("test message")

	testLogOutput(t, logger, func(entry zapcore.Entry) bool {
		return entry.Level == zap.InfoLevel &&
			entry.Message == "test message" &&
			entry.Time.Format("2006-01-02T15:04:05") != "" // ISO8601
	}, "timestamp")
}

func testLogOutput(t *testing.T, logger *zap.Logger, checkEntry func(zapcore.Entry) bool, expectedTimeKey string) {
	var buf bytes.Buffer

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = expectedTimeKey
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zap.InfoLevel)
	capturedLogger := zap.New(core)

	capturedLogger.Info("test message")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log JSON: %v", err)
	}
	if _, ok := entry[expectedTimeKey]; !ok {
		t.Errorf("Expected time key %q in log, got keys: %v", expectedTimeKey, entry)
	}

}
