package main

import "testing"

func TestUI(t *testing.T) {
	want := "test_OK"
	if got := forTesting(); got != want {
		t.Errorf("forTesting() = %q, want %q", got, want)
	}
}
