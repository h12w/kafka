package broker

import (
	"testing"
)

func TestErrorEqual(t *testing.T) {
	t.Parallel()
	var err error
	err = ErrOffsetOutOfRange
	if err != ErrOffsetOutOfRange {
		t.Fatal("error should be equal")
	}
}
