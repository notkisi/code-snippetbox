package assert

import "testing"

func Equal[T comparable](t *testing.T, actual T, expected T) {
	// marks this method as helper function; this skips this method in logs
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v, expected: %v", actual, expected)
	}
}
