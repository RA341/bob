package vm

import (
	"cmp"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed)

// VariadicArgCount indicates a variadic fn
const VariadicArgCount = -1

type FnDef struct {
	ArgCount int
	Fn       func(args ...string) error
}

func NewFnDev(
	argCount int,
	Fn func(args ...string) error,
) FnDef {
	return FnDef{
		Fn:       Fn,
		ArgCount: argCount,
	}
}

type VM struct {
	program []Ins
	stack   Stack[Value]

	callStack Stack[int]
	labels    map[string]int

	Builtins map[string]FnDef
	Vars     map[string]Value
	pc       int
}

func (vm *VM) Start(ins []Ins, funcs map[string]FnDef) {
	vm.program = ins
	vm.pc = 0
	vm.labels = make(map[string]int)

	vm.Vars = make(map[string]Value)

	vm.Builtins = funcs
	if funcs == nil {
		vm.Builtins = make(map[string]FnDef)
	}

	vm.Run()
}

func (vm *VM) Run() {
	for i, ins := range vm.program {
		if ins.OpType == LABEL {
			vm.labels[ins.Value.Raw] = i
		}
	}

	for range vm.program {
		if vm.pc >= len(vm.program) {
			log.Fatal(
				"vm attempted to get a instruction that is out of bounds",
				"Max allowed",
				len(vm.program),
				"Attempted",
				vm.pc,
			)
		}

		ins := vm.program[vm.pc]
		vm.pc++

		if vm.executeInstruction(&ins) {
			return
		}
	}
}

func (vm *VM) executeInstruction(ins *Ins) bool {
	switch ins.OpType {
	case PUSH:
		vm.stack.Push(ins.Value)
	case STORE:
		vm.Store()
	case CALL:
		vm.Call(ins)
	case LOAD:
		val, ok := vm.Vars[ins.Value.Raw]
		if !ok {
			log.Fatal("Could not find variable ", ins.Value.Raw)
		}
		vm.stack.Push(val)
	case ADD:
		vm.Add()
	case LABEL:
		// no op only indicate body of the next instruction
	case JZ:
		val := vm.mustPopBool()
		if val {
			jmpIdx, ok := vm.labels[ins.Value.Raw]
			if !ok {
				log.Fatalf("label '%q' does not exist", ins.Value.Raw)
			}
			vm.pc = jmpIdx
		}
	case JMP:
		vm.pc = vm.labels[ins.Value.Raw]
	case JNZ:
		val := vm.mustPopBool()
		if !val {
			jmpIdx, ok := vm.labels[ins.Value.Raw]
			if !ok {
				log.Fatalf("label '%q' does not exist", ins.Value.Raw)
			}
			vm.pc = jmpIdx
		}
	case RET:
		retPc := vm.callStack.MustPop()
		vm.pc = retPc

	case EQ:
		cond1 := vm.stack.MustPop()
		cond2 := vm.stack.MustPop()

		res := cond1.Raw == cond2.Raw
		vm.stack.Push(BoolVal(res))
	case GT:
		vm.Compare(true)
	case LT:
		vm.Compare(false)

	case NOT:
		v := vm.mustPopBool()
		vm.stack.Push(BoolVal(!v))
	case AND:
		cond1 := vm.mustPopBool()
		cond2 := vm.mustPopBool()

		res := cond1 && cond2
		vm.stack.Push(BoolVal(res))
	case OR:
		cond1 := vm.mustPopBool()
		cond2 := vm.mustPopBool()

		res := cond1 || cond2
		vm.stack.Push(BoolVal(res))
	case HALT:
		log.Println("VM exited successfully")
		return true
	default:
		log.Fatal(red.Sprintf("Unsupported op type: %q", ins.OpType))
	}
	return false
}

func (vm *VM) mustPopBool() bool {
	popBool, err := vm.popBool()
	if err != nil {
		log.Fatal(err)
	}
	return popBool
}

func (vm *VM) popBool() (bool, error) {
	val, err := vm.stack.Pop()
	if err != nil {
		return false, err
	}

	switch val.Type {
	case VTString:
		return val.Raw != "", nil
	case VTInt, VTFloat:
		return val.Raw != "0", nil
	case VTBool:
		bVal, err := strconv.ParseBool(val.Raw)
		if err != nil {
			return false, fmt.Errorf("could not parse bool %q: %w", val.Raw, err)
		}
		return bVal, nil
	default:
		return false, fmt.Errorf("unknown type: %v: %s", val.Type, val.Raw)
	}
}

func (vm *VM) popAssert(vt ValueType) Value {
	val := vm.stack.MustPop()
	if val.Type != vt {
		log.Fatalf(
			"expected %s type got [%s]: %s",
			vt.String(),
			val.Type,
			val.Raw,
		)
	}
	return val
}

func (vm *VM) Store() {
	varName := vm.stack.MustPop()
	varVal := vm.stack.MustPop()
	vm.Vars[varName.Raw] = varVal
}

func (vm *VM) Call(ins *Ins) {
	fnName := vm.stack.MustPop()

	fnDef, ok := vm.Builtins[fnName.Raw]
	if !ok {
		log.Fatal(red.Sprintf("no function found with name name: %q", fnName))
	}

	argCount := fnDef.ArgCount

	var err error
	if argCount == VariadicArgCount {
		varCountStr := ins.Value.Raw
		argCount, err = strconv.Atoi(varCountStr)
		if err != nil {
			log.Fatal(
				red.Sprintf(
					"invalid variadic count: %s\n"+
						"err: %v\n"+
						"This function is variadic and a explicit variable count IS REQUIRED",
					varCountStr,
					err,
				),
			)
		}
	}

	args := make([]string, argCount)
	for i := range argCount {
		argStr := vm.stack.MustPop()
		args[i] = argStr.Raw
	}

	err = fnDef.Fn(args...)
	if err != nil {
		log.Fatal(red.Sprintf("error calling function %q: %s", fnName, err))
	}
}

func (vm *VM) Add() {
	v1 := vm.stack.MustPop()
	v2 := vm.stack.MustPop()

	switch {
	case v1.Type == VTInt && v2.Type == VTInt:
		i1, _ := strconv.Atoi(v1.Raw)
		i2, _ := strconv.Atoi(v2.Raw)
		vm.stack.Push(IntVal(i1 + i2))

	case v1.Type == VTFloat || v2.Type == VTFloat:
		f1, _ := strconv.ParseFloat(v1.Raw, 64)
		f2, _ := strconv.ParseFloat(v2.Raw, 64)
		vm.stack.Push(FloatVal(f1 + f2))

	case v1.Type == VTString && v2.Type == VTString:
		vm.stack.Push(StrVal(v1.Raw + v2.Raw))

	default:
		log.Fatal(red.Sprintf("type mismatch: cannot add %v and %v", v1.Type, v2.Type))
	}
}

func (vm *VM) Compare(shouldBeGreater bool) {
	v1 := vm.stack.MustPop()
	v2 := vm.stack.MustPop()

	var res int

	switch {
	case v1.Type == VTInt && v2.Type == VTInt:
		i1, _ := strconv.Atoi(v1.Raw)
		i2, _ := strconv.Atoi(v2.Raw)

		res = cmp.Compare(i1, i2)

	case v1.Type == VTFloat || v2.Type == VTFloat:
		f1, _ := strconv.ParseFloat(v1.Raw, 64)
		f2, _ := strconv.ParseFloat(v2.Raw, 64)

		res = cmp.Compare(f1, f2)
	case v1.Type == VTString && v2.Type == VTString ||
		v1.Type == VTBool && v2.Type == VTBool:
		res = cmp.Compare(v1.Raw, v2.Raw)
	default:
		log.Fatal(red.Sprintf("type mismatch: cannot compare %v and %v", v1.Type, v2.Type))
	}

	if shouldBeGreater {
		vm.stack.Push(BoolVal(res == 1))
	} else {
		vm.stack.Push(BoolVal(res == -1))
	}

}
