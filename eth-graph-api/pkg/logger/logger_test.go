package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestInitialize(t *testing.T) {
	err := Initialize()
	if err != nil {
		t.Errorf("Failed to initialize logger: %v", err)
	}
	if Log == nil {
		t.Error("Expected logger to be initialized, but got nil")
	}
}

func TestInfo(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)

	Log = zap.New(core).Sugar()

	message := "test info log"
	Log.Infow(message, "key", "value")

	if logs.Len() != 1 || logs.All()[0].Message != message || logs.All()[0].Context[0].String != "value" {
		t.Error("Unexpected log message or context")
	}
}

func TestError(t *testing.T) {
	core, logs := observer.New(zap.ErrorLevel)

	Log = zap.New(core).Sugar()

	message := "test error log"
	Log.Errorw(message, "key", "value")

	if logs.Len() != 1 || logs.All()[0].Message != message || logs.All()[0].Context[0].String != "value" {
		t.Error("Unexpected log message or context")
	}
}
