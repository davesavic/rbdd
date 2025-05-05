package app

import (
	"os"
	"testing"
)

func TestIExecuteCommand(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommand("echo test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iExecuteCommand("thiscommanddoesnotexist")
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}

	err = apiTest.iExecuteCommand("")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

func TestIExecuteCommandInDirectory(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	tempDir, err := os.MkdirTemp("", "apitest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = apiTest.iExecuteCommandInDirectory("pwd", tempDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	apiTest.store["dir"] = tempDir
	err = apiTest.iExecuteCommandInDirectory("pwd", "${dir}")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestIExecuteCommandWithTimeout(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommandWithTimeout("echo test", 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.iExecuteCommandWithTimeout("sleep 2", 1)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestTheCommandOutputShouldMatch(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommand("echo test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theCommandOutputShouldMatch("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theCommandOutputShouldMatch("not test")
	if err == nil {
		t.Error("Expected error for non-matching output, got nil")
	}
}

func TestTheCommandOutputShouldContain(t *testing.T) {
	apiTest := NewAPITest("https://example.com")

	err := apiTest.iExecuteCommand("echo test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theCommandOutputShouldContain("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = apiTest.theCommandOutputShouldContain("not test")
	if err == nil {
		t.Error("Expected error for non-matching output, got nil")
	}
}
