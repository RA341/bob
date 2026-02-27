package vm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVM_State(t *testing.T) {
	state := NewVmState[string]()

	value1 := "testV 1"
	state.Push(value1)

	value2 := "testV 2"
	state.Push(value2)

	value3 := "testV 3"
	state.Push(value3)

	value4 := "test 4"
	state.Push(value4)

	require.Equal(t, value4, state.MustPop())
	require.Equal(t, value3, state.MustPop())

	value5 := "testV 5"
	state.Push(value5)
	require.Equal(t, value5, state.MustPop())

	require.Equal(t, value2, state.MustPop())
	require.Equal(t, value1, state.MustPop())

	require.Equal(t, state.Len(), 0)
}
