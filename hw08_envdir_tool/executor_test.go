package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Expected return code", func(t *testing.T) {
		status := RunCmd([]string{"bash", "-c", "exit 10"}, Environment{})
		require.Equal(t, 10, status)
	})

	t.Run("Zero code in normal way", func(t *testing.T) {
		status := RunCmd([]string{"echo", "arg1=1", "arg2=2"}, Environment{})
		require.Equal(t, 0, status)
	})
}
