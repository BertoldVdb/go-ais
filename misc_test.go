package ais

import (
	"testing"
)

func TestAssert(t *testing.T) {
	assert(true, "No error")

	defer func() {
		recover()
	}()

	assert(false, "Error")
	t.Errorf("Did not panic")
}
