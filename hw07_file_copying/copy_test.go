package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy file with unknown size", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("invalid offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/", math.MaxInt64, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})
}
