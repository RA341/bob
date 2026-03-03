package flag

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

	sd := Parser{
		Flags: []*Flag{
			Bool("version", ""),
			StrFlag("value", ""),
			Bool("bool", ""),
			IntFlag("wow", ""),
		},
	}

	err := sd.Parse(in)
	require.NoError(t, err)

	require.Equal(t, sd.Flags[0].Val, false)
	require.Equal(t, sd.Flags[0].IsSet, false)

	require.Equal(t, sd.Flags[1].Val, "someValue")
	require.Equal(t, sd.Flags[1].IsSet, true)

	require.Equal(t, sd.Flags[2].Val, true)
	require.Equal(t, sd.Flags[2].IsSet, true)

	require.Equal(t, sd.Flags[3].Val, 98234)
	require.Equal(t, sd.Flags[3].IsSet, true)

	require.Equal(t, in[sd.Current], in[5])
}

func Test_ParseFlags_Err(t *testing.T) {
	sd := Parser{
		Flags: []*Flag{
			IntFlag("wow", ""),
			Bool("version", ""),
			StrFlag("value", ""),
			Bool("bool", ""),
		},
	}

	in := []string{
		"--value",
	}
	err := sd.Parse(in)
	require.Error(t, err)
	t.Log(err)

	sd = Parser{
		Flags: []*Flag{
			IntFlag("wow", ""),
			Bool("version", ""),
			StrFlag("value", ""),
			Bool("bool", ""),
		},
	}
	in = []string{
		"--wow",
		"asdads",
	}

	err = sd.Parse(in)
	require.Error(t, err)
	t.Log(err)
}
