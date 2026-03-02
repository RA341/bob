package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseFlags(t *testing.T) {
	in := []string{
		"--wow", "98234",
		"--bool",
		"--value", "someValue",
		"subcommand",
	}

	sd := FlagParser{
		flags: []*Flag{
			BoolFlag("version", ""),
			StrFlag("value", ""),
			BoolFlag("bool", ""),
			IntFlag("wow", ""),
		},
	}

	err := sd.parse(in)
	require.NoError(t, err)

	require.Equal(t, sd.flags[0].val, false)
	require.Equal(t, sd.flags[0].isSet, false)

	require.Equal(t, sd.flags[1].val, "someValue")
	require.Equal(t, sd.flags[1].isSet, true)

	require.Equal(t, sd.flags[2].val, true)
	require.Equal(t, sd.flags[2].isSet, true)

	require.Equal(t, sd.flags[3].val, 98234)
	require.Equal(t, sd.flags[3].isSet, true)

	require.Equal(t, in[sd.current], in[5])
}

func Test_ParseFlags_Err(t *testing.T) {
	sd := FlagParser{
		flags: []*Flag{
			IntFlag("wow", ""),
			BoolFlag("version", ""),
			StrFlag("value", ""),
			BoolFlag("bool", ""),
		},
	}

	in := []string{
		"--value",
	}
	err := sd.parse(in)
	require.Error(t, err)
	t.Log(err)

	sd = FlagParser{
		flags: []*Flag{
			IntFlag("wow", ""),
			BoolFlag("version", ""),
			StrFlag("value", ""),
			BoolFlag("bool", ""),
		},
	}
	in = []string{
		"--wow",
		"asdads",
	}
	err = sd.parse(in)
	require.Error(t, err)
	t.Log(err)
}
