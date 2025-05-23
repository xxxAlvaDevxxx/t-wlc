package twlc

import (
	"os"
	"testing"
)

func TestDefaultTwlc(t *testing.T) {
	twlc := DefaultTwlc()
	if twlc.SaveInLogFile != true {
		t.Errorf("Expected SaveInLogFile to be true, got %v", twlc.SaveInLogFile)
	}
	if twlc.ShowInConsole != true {
		t.Errorf("Expected ShowInConsole to be true, got %v", twlc.ShowInConsole)
	}
	if twlc.ColorMessages != true {
		t.Errorf("Expected ColorMessages to be true, got %v", twlc.ColorMessages)
	}
	if twlc.LogDir == "" {
		t.Errorf("Expected LogDir to be non-empty, got %s", twlc.LogDir)
	}
	if twlc.WithTime != false {
		t.Errorf("Expected WithTime to be true, got %v", twlc.WithTime)
	}

	// Test createLogDir
	_, err := os.Stat(twlc.LogDir)
	if os.IsNotExist(err) {
		t.Errorf("Expected log directory to exist, got %v", err)
	}
}

func TestNewTwlc(t *testing.T) {
	twlc := NewTwlc(true, true, true, true, "./test_logs/")
	if twlc.SaveInLogFile != true {
		t.Errorf("Expected SaveInLogFile to be true, got %v", twlc.SaveInLogFile)
	}
	if twlc.ShowInConsole != true {
		t.Errorf("Expected ShowInConsole to be true, got %v", twlc.ShowInConsole)
	}
	if twlc.LogDir != "./test_logs/" {
		t.Errorf("Expected LogDir to be /tmp/test_logs/, got %s", twlc.LogDir)
	}
	if twlc.WithTime != true {
		t.Errorf("Expected WithTime to be true, got %v", twlc.WithTime)
	}
	if twlc.ColorMessages != true {
		t.Errorf("Expected ColorMessages to be true, got %v", twlc.ColorMessages)
	}

	// Check if the log directory was created
	_, err := os.Stat(twlc.LogDir)
	if os.IsNotExist(err) {
		t.Errorf("Expected log directory to exist, got %v", err)
	}

	// Clean up
	err = os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}

}

func TestWriteLog(t *testing.T) {
	twlc := NewTwlc(true, true, true, true, "./test_logs/")
	message := "Test message"
	twlc.WriteLog(Info, message)

	// Check if the log file was created
	logFilePath := twlc.LogDir + "log.txt"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Errorf("Expected log file to exist, got %v", err)
	}

	// Clean up
	err := os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}
}

func TestWriteConstants(t *testing.T) {
	twlc := NewTwlc(true, true, true, true, "./test_logs/")
	twlc.WriteInfo(Info)
	twlc.WriteSuccess(Success)
	twlc.WriteError(Error)
	twlc.WriteWarning(Warning)
	twlc.WriteDebug(Debug)
	twlc.WriteTrace(Trace)

	// Check if the log file was created
	logFilePath := twlc.LogDir + "log.txt"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Errorf("Expected log file to exist, got %v", err)
	}
	// Clean up
	err := os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}
}

func TestSetColor(t *testing.T) {
	cases := []struct {
		messageType string
		expected    string
	}{
		{Info, "\033[34mTest message\033[0m"},
		{Success, "\033[32mTest message\033[0m"},
		{Warning, "\033[33mTest message\033[0m"},
		{Error, "\033[31mTest message\033[0m"},
		{Debug, "\033[35mTest message\033[0m"},
		{Trace, "\033[36mTest message\033[0m"},
	}

	twlc := NewTwlc(true, true, true, true, "./test_logs/")

	for _, c := range cases {
		result := twlc.setColor(c.messageType, "Test message")
		if result != c.expected {
			t.Errorf("Expected %s, got %s", c.expected, result)
		}

	}

	// Test with an unknown message type
	unknownMessageType := "Unknown"
	unknownMessage := twlc.setColor(unknownMessageType, "Test message")
	if unknownMessage != "Test message" {
		t.Errorf("Expected 'Test message', got %s", unknownMessage)
	}

	// Clean up
	err := os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}
}

func TestStructToString(t *testing.T) {
	// Create structs for testing
	type Animal struct {
		Name string
		Age  int
	}

	cases := []struct {
		input    interface{}
		expected string
		simple   bool
	}{
		{Animal{"Dog", 5}, "{Name:Dog Age:5}", true},
		{Animal{"Cat", 3}, "{Name:Cat Age:3}", true},
		{Animal{"Dog", 5}, `twlc.Animal{Name:"Dog", Age:5}`, false},
	}
	twlc := NewTwlc(true, true, true, true, "./test_logs/")
	for _, c := range cases {
		result := twlc.StructToString(c.input, c.simple)
		if result != c.expected {
			t.Errorf("Expected %s, got %s", c.expected, result)
		}
	}

	// Clean up
	err := os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}
}

func TestStructToJson(t *testing.T) {
	// Create structs for testing
	type Animal struct {
		Name string
		Age  int
	}

	cases := []struct {
		input    interface{}
		expected string
	}{
		{Animal{"Dog", 5}, `{
    "Name": "Dog",
    "Age": 5
}`},
		{Animal{"Cat", 3}, `{
    "Name": "Cat",
    "Age": 3
}`},
	}

	twlc := NewTwlc(true, true, true, true, "./test_logs/")
	for _, c := range cases {
		result, err := twlc.StructToJson(c.input)
		if err != nil {
			t.Errorf("Failed to convert struct to JSON: %v", err)
		}
		if result != c.expected {
			t.Errorf("Expected %s, got %s", c.expected, result)
		}
	}

	// Clean up
	err := os.RemoveAll(twlc.LogDir)
	if err != nil {
		t.Errorf("Failed to remove log directory: %v", err)
	}
}
