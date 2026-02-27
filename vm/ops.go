package vm

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go run github.com/dmarkham/enumer@latest -type=OpType -output=gen_enum_op.go
type OpType int

const (
	PUSH OpType = iota
	POP

	// STORE assigns a var
	STORE
	// CALL calls a builtin function
	CALL
	// LOAD resolves a var from the stack
	LOAD
	ADD

	HALT
	LABEL
	RET

	// JMP unconditional jump
	JMP
	// JZ Jump if zero (true)
	JZ
	// JNZ Jump if not zero (false)
	JNZ

	NOT
	OR
	AND

	EQ
	// GT last value popped is greater than first
	GT
	// LT last value popped is less than first
	LT
)

//go:generate go run github.com/dmarkham/enumer@latest -type=ValueType -output=gen_enum_value.go
type ValueType int

const (
	VTString ValueType = iota
	VTInt
	VTFloat
	VTBool
)

func IntVal(n int) Value {
	return Value{VTInt, strconv.Itoa(n)}
}

func FloatVal(f float64) Value {
	return Value{VTFloat, strconv.FormatFloat(f, 'f', -1, 64)}
}

func StrVal(s string) Value {
	return Value{VTString, s}
}

func BoolVal(val bool) Value {
	return Value{VTBool, strconv.FormatBool(val)}
}

type Value struct {
	Type ValueType
	Raw  string // keep raw for display/storage
}

func (i Value) String() string {
	return fmt.Sprintf(
		"[%s] %s",
		strings.TrimSuffix(strings.TrimPrefix(i.Type.String(), "VT"), "ing"),
		i.Raw,
	)
}

// Ins Reprints a single instruction
type Ins struct {
	OpType OpType
	Value  Value
}

func (i Ins) String() string {
	return fmt.Sprintf("[%s] %v", i.OpType.String(), i.Value)
}

// O loads an op with only an op (represents a no value op)
func O(op OpType) Ins {
	return Ins{OpType: op}
}

// OV inits an instruction with op and val
func OV(op OpType, val Value) Ins {
	return Ins{
		OpType: op,
		Value:  val,
	}
}

func OVInt(op OpType, val int) Ins {
	return Ins{OpType: op, Value: IntVal(val)}
}

func OVStr(op OpType, val string) Ins {
	return Ins{OpType: op, Value: StrVal(val)}
}

func OVFloat(op OpType, val float64) Ins {
	return Ins{OpType: op, Value: FloatVal(val)}
}

func OVBool(op OpType, val bool) Ins {
	return Ins{OpType: op, Value: BoolVal(val)}
}

func NewValue(valType ValueType, rawVal string) Value {
	return Value{
		Type: valType,
		Raw:  rawVal,
	}
}
