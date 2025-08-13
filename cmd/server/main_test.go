package main

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "Environment variable exists",
			envKey:       "TEST_ENV_VAR",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "Environment variable does not exist",
			envKey:       "NON_EXISTENT_VAR",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Environment variable is empty",
			envKey:       "EMPTY_VAR",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Default value is empty",
			envKey:       "ANOTHER_NON_EXISTENT_VAR",
			envValue:     "",
			defaultValue: "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := getEnvOrDefault(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnvOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnvOrDefault_RealEnvVar(t *testing.T) {
	// Test with PATH environment variable (should exist on all systems)
	result := getEnvOrDefault("PATH", "default_path")

	// PATH should exist and not be empty on any system
	if result == "default_path" {
		t.Error("Expected PATH environment variable to exist")
	}

	if result == "" {
		t.Error("Expected PATH environment variable to not be empty")
	}
}

func TestGetEnvOrDefault_WithSpaces(t *testing.T) {
	// Test with values containing spaces
	envKey := "TEST_SPACES_VAR"
	envValue := "value with spaces"
	defaultValue := "default with spaces"

	os.Setenv(envKey, envValue)
	defer os.Unsetenv(envKey)

	result := getEnvOrDefault(envKey, defaultValue)
	if result != envValue {
		t.Errorf("getEnvOrDefault() = %v, want %v", result, envValue)
	}
}

func TestGetEnvOrDefault_WithSpecialCharacters(t *testing.T) {
	// Test with special characters
	envKey := "TEST_SPECIAL_VAR"
	envValue := "value!@#$%^&*()"
	defaultValue := "default!@#"

	os.Setenv(envKey, envValue)
	defer os.Unsetenv(envKey)

	result := getEnvOrDefault(envKey, defaultValue)
	if result != envValue {
		t.Errorf("getEnvOrDefault() = %v, want %v", result, envValue)
	}
}
