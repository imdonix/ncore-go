package ncore

import (
	"testing"
)

func TestSizeEqual(t *testing.T) {
	tests := []struct {
		s1 string
		s2 string
	}{
		{"1024 MiB", "1 GiB"},
		{"10 MiB", "10 MiB"},
		{"2048 KiB", "2 MiB"},
	}

	for _, tc := range tests {
		sz1, _ := NewSize(tc.s1)
		sz2, _ := NewSize(tc.s2)
		if sz1.Bytes() != sz2.Bytes() {
			t.Errorf("Expected %s to be equal to %s", tc.s1, tc.s2)
		}
	}
}

func TestSizeAdd(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected string
	}{
		{"1024 MiB", "1 GiB", "2.00 GiB"},
		{"10 MiB", "11 MiB", "21.00 MiB"},
		{"2048 KiB", "2 MiB", "4.00 MiB"},
	}

	for _, tc := range tests {
		sz1, _ := NewSize(tc.s1)
		sz2, _ := NewSize(tc.s2)
		result := sz1.Add(sz2)
		if result.String() != tc.expected {
			t.Errorf("Expected %s + %s to be %s, got %s", tc.s1, tc.s2, tc.expected, result.String())
		}
	}
}
