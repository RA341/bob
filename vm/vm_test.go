package vm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmp(t *testing.T) {
	vm := new(VM)
	input := []Ins{
		OVStr(PUSH, "test"),
		OVStr(PUSH, "@value"),
		O(LT),
	}
	run(vm, input)

	res, err := vm.popBool()
	require.NoError(t, err)
	require.Equal(t, true, res)

	input = []Ins{
		OVInt(PUSH, 10),
		OVInt(PUSH, 20),
		O(GT),
	}
	run(vm, input)

	res, err = vm.popBool()
	require.NoError(t, err)
	require.Equal(t, true, res)

	input = []Ins{
		OVInt(PUSH, 50),
		OVInt(PUSH, 20),
		O(GT),
	}
	run(vm, input)

	res, err = vm.popBool()
	require.NoError(t, err)
	require.NotEqual(t, true, res)

}

func TestVm(t *testing.T) {
	vm := new(VM)

	input := []Ins{
		OVStr(PUSH, "test"),
		OVStr(PUSH, "@value"),
		O(STORE),

		OVStr(LOAD, "@value"),
		OVStr(PUSH, "This value was inserted using formatting"),
		OVStr(PUSH, "%s, %s\n"),
		OVStr(PUSH, "printf"),
		OVInt(CALL, 3),

		OVInt(PUSH, 1),
		OVInt(PUSH, 30),
		O(ADD),

		OVStr(PUSH, "test"),
		OVStr(LOAD, "@value"),
		O(ADD),
		OVStr(LOAD, "@value"),
		O(ADD),

		OVStr(PUSH, "@value"),
		O(STORE),

		OVStr(LOAD, "@value"),
		OVStr(PUSH, "print"),
		OVInt(CALL, 1),
		O(HALT),
	}

	run(vm, input)
}

func TestVm2(t *testing.T) {
	vm := new(VM)

	input := []Ins{
		OVStr(PUSH, "someValue"),
		OVStr(PUSH, "@value"),
		OVStr(PUSH, "someValue"),
		OVStr(PUSH, "@value2"),
		O(EQ),

		OVStr(PUSH, "%s %s %s\n"),
		OVStr(PUSH, "@value2"),

		O(HALT),
		OVStr(PUSH, "@value"),

		OVStr(PUSH, "printf"),
		OVInt(CALL, 3),
	}

	run(vm, input)
}

func run(vm *VM, input []Ins) {
	vm.Start(input, DefaultFns)

	fmt.Println()
	fmt.Println("=================================")

	fmt.Println("Vars", vm.Vars)
	fmt.Println("=================================")
	fmt.Println("Stack", vm.stack)
	fmt.Println("=================================")
	fmt.Println("Program", vm.program)

	fmt.Println()
	fmt.Println("=================================")
}
