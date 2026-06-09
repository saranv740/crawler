package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:          "remove https",
			input:         "https://www.boot.dev/blog/path",
			expected:      "www.boot.dev/blog/path",
			expectedError: nil,
		},
		{
			name:          "remove trailing slash",
			input:         "https://www.boot.dev/blog/path/",
			expected:      "www.boot.dev/blog/path",
			expectedError: nil,
		},
		{
			name:          "remove http",
			input:         "http://www.boot.dev/blog/path",
			expected:      "www.boot.dev/blog/path",
			expectedError: nil,
		},
		{
			name:          "remove capitalization",
			input:         "https://www.boot.dev/blog/PATH",
			expected:      "www.boot.dev/blog/path",
			expectedError: nil,
		},
		{
			name:          "normalize slashes",
			input:         `https://www.boot.dev//blog\PATH\\hello//`,
			expected:      "www.boot.dev/blog/path/hello",
			expectedError: nil,
		},
		{
			name:          "invalid url",
			input:         `://www.boot.dev/`,
			expected:      "",
			expectedError: ErrInvalidURL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := normalizeURL(tc.input)
			if err != tc.expectedError {
				t.Fatalf("test %q failed. expected %v, got %v", tc.name, tc.expectedError, err)
			}

			if out != tc.expected {
				t.Fatalf("test %q failed. expected %q, got %q", tc.name, tc.expected, out)
			}
		})
	}
}
