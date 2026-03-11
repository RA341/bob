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
	Stack   Stack[Value]

	callStack Stack[int]
	labels    map[string]int

	Builtins map[string]FnDef
	Vars     map[string]Value
	pc       int

	errs []error
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

	vm.run()
}

func (vm *VM) run() {
	for i, ins := range vm.program {
		if ins.OpType == LABEL {
			vm.labels[ins.Value.Raw] = i
		}
	}

	for range vm.program {
		if vm.pc >= len(vm.program) {
			vm.errs = append(vm.errs, fmt.Errorf("vm stack overflow"))

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
		vm.Stack.Push(ins.Value)
	case STORE:
		vm.store()
	case CALL:
		vm.call()
	case LOAD:
		val, ok := vm.Vars[ins.Value.Raw]
		if !ok {
			log.Fatal("Could not find variable ", ins.Value.Raw)
		}
		vm.Stack.Push(val)
	case ADD:
		vm.add()
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
		cond1 := vm.Stack.MustPop()
		cond2 := vm.Stack.MustPop()

		res := cond1.Raw == cond2.Raw
		vm.Stack.Push(BoolVal(res))
	case GT:
		vm.compare(true)
	case LT:
		vm.compare(false)

	case NOT:
		v := vm.mustPopBool()
		vm.Stack.Push(BoolVal(!v))
	case AND:
		cond1 := vm.mustPopBool()
		cond2 := vm.mustPopBool()

		res := cond1 && cond2
		vm.Stack.Push(BoolVal(res))
	case OR:
		cond1 := vm.mustPopBool()
		cond2 := vm.mustPopBool()

		res := cond1 || cond2
		vm.Stack.Push(BoolVal(res))
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
	val, err := vm.Stack.Pop()
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
	val := vm.Stack.MustPop()
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

func (vm *VM) store() {
	varName := vm.Stack.MustPop()
	varVal := vm.Stack.MustPop()
	vm.Vars[varName.Raw] = varVal
}

func (vm *VM) call() {
	fnName := vm.Stack.MustPop()

	fnDef, ok := vm.Builtins[fnName.Raw]
	if !ok {
		log.Fatal(red.Sprintf("no function found with name name: %q", fnName))
	}

	varCountStr := vm.Stack.MustPop()
	argCount, err := strconv.Atoi(varCountStr.Raw)
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

	args := make([]string, argCount)
	for i := range argCount {
		argStr := vm.Stack.MustPop()
		args[i] = argStr.Raw
	}

	err = fnDef.Fn(args...)
	if err != nil {
		log.Fatal(red.Sprintf("error calling function %q: %s", fnName, err))
	}
}

func (vm *VM) add() {
	v1 := vm.Stack.MustPop()
	v2 := vm.Stack.MustPop()

	switch {
	case v1.Type == VTInt && v2.Type == VTInt:
		i1, _ := strconv.Atoi(v1.Raw)
		i2, _ := strconv.Atoi(v2.Raw)
		vm.Stack.Push(IntVal(i1 + i2))

	case v1.Type == VTFloat || v2.Type == VTFloat:
		f1, _ := strconv.ParseFloat(v1.Raw, 64)
		f2, _ := strconv.ParseFloat(v2.Raw, 64)
		vm.Stack.Push(FloatVal(f1 + f2))

	case v1.Type == VTString && v2.Type == VTString:
		vm.Stack.Push(StrVal(v1.Raw + v2.Raw))

	default:
		log.Fatal(red.Sprintf("type mismatch: cannot add %v and %v", v1.Type, v2.Type))
	}
}

func (vm *VM) compare(shouldBeGreater bool) {
	v1 := vm.Stack.MustPop()
	v2 := vm.Stack.MustPop()

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
		vm.Stack.Push(BoolVal(res == 1))
	} else {
		vm.Stack.Push(BoolVal(res == -1))
	}

}
