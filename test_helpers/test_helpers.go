package test_helpers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

/**
 * round float to precision - should be used when comparing floats
 * as precision errors can occur when doing arithmetic operations
 */
func RoundFloat(val, precision float64) float64 {
	return float64(int64(val/precision+0.5)) * precision
}

func Assert(t *testing.T, actual bool, message string, tokens ...interface{}) {
	if !actual {
		t.Errorf(message, tokens...)
	}
}

func AssertNotNil(t *testing.T, actual interface{}, optionalMessageAndTokens ...interface{}) {
	if actual == nil {
		failMessage := withUserMessage("Expected actual to be not nil", optionalMessageAndTokens...)
		t.Errorf(failMessage)
	}
}

func AssertNil(t *testing.T, actual interface{}, optionalMessageAndTokens ...interface{}) {
	if actual != nil {
		failMessage := withUserMessage("Expected %v (type: %T) to be nil", optionalMessageAndTokens...)
		t.Errorf(failMessage, actual, actual)
	}
}

func AssertTimeEqual(t *testing.T, expected, actual time.Time, optionalMessageAndTokens ...interface{}) {
	if !expected.Equal(actual) {
		failMessage := withUserMessage("Expected %v to be %v", optionalMessageAndTokens...)
		t.Errorf(failMessage, expected, actual)
	}
}

func AssertDeepEqual(t *testing.T, expected, actual interface{}, optionalMessageAndTokens ...interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		failMessage := withUserMessage("Expected ==:\n  expected: %v\n  actual:   %v", optionalMessageAndTokens...)
		t.Errorf(failMessage, expected, actual)
	}
}

func AssertEqual(t *testing.T, expected, actual interface{}, optionalMessageAndTokens ...interface{}) {
	if expected != actual {
		failMessage := withUserMessage("Expected ==:\n  expected: %v\n  actual:   %v", optionalMessageAndTokens...)
		t.Errorf(failMessage, expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected, actual interface{}, optionalMessageAndTokens ...interface{}) {
	if expected == actual {
		failMessage := withUserMessage("Expected !=:\n  expected: %v\n  actual:   %v", optionalMessageAndTokens...)
		t.Errorf(failMessage, expected, actual)
	}
}

func withUserMessage(defaultMessage string, stringAndMaybeTokens ...interface{}) string {
	if len(stringAndMaybeTokens) == 0 {
		return defaultMessage
	}

	if additionalMessage, isStr := stringAndMaybeTokens[0].(string); isStr {
		if additionalMessage != "" {
			tokens := stringAndMaybeTokens[1:]
			if len(tokens) > 0 {
				additionalMessage = fmt.Sprintf(additionalMessage, tokens...)
			}
			defaultMessage = fmt.Sprintf("%s;\n\t\t%s", additionalMessage, defaultMessage)
		}
	}

	return defaultMessage
}

func contains(arr []string, candidate string) bool {
	for _, v := range arr {
		if v == candidate {
			return true
		}
	}

	return false
}

func ExpectsError(t *testing.T, expectedMessage string, err error) {
	if err == nil {
		t.Errorf("Expected error [%s], but received none", expectedMessage)
		return
	}

	if expectedMessage != "" && !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("Failed to receive correct error message\n  expected: [%s]\n    actual: [%v]", expectedMessage, err)
	}
}

func AssertContains(t *testing.T, arr []string, candidate string) {
	if !contains(arr, candidate) {
		t.Errorf("Expected %v to contain %v, but didn't.", arr, candidate)
	}
}
