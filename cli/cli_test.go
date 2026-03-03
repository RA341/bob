package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelp(t *testing.T) {
	ap := NewApp()

	args := []string{
		"bob",
		"--help",
	}

	err := ap.Run(args)
	require.NoError(t, err)
}
