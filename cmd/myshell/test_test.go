package main

import "testing"

func TestCase1(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "test", expected: "test"},
		{input: "test    test", expected: "test test"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			output := resolveArguments(tt.input)
			if output != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, output)
			}
		})
	}

}

func TestTest(t *testing.T) {
	_ = resolveArguments("test")
}
