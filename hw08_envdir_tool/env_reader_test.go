package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Incorrect filename with =", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(".", "testdir_")
		require.NoError(t, err)
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				require.FailNow(t, err.Error())
			}
		}(tempDir)

		_, err = os.CreateTemp(tempDir, "test=ignored")
		require.NoError(t, err)

		env, err := ReadDir(tempDir)
		require.Equal(t, ErrorIncorrectFileName, err)
		require.Empty(t, env)
	})
	t.Run("Collect data", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		expected := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		require.Equal(t, expected, env)
		require.NoError(t, err)
	})
}
