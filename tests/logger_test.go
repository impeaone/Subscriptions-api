package tests

import "testing"

func TestLogger(t *testing.T) {
	logger := NewTestLog("")
	logger.Info("Info", "test")
	logger.Warning("Warning", "test")
	logger.Error("Error", "test")
}
