package vm

import (
	"errors"
	"fmt"
	"log"
)

type Stack[T any] struct {
	arr []T
}

func NewVmState[T any]() Stack[T] {
	return Stack[T]{}
}

func (s *Stack[T]) Len() int {
	return len(s.arr)
}

func (s *Stack[T]) Push(value T) {
	s.arr = append(s.arr, value)
}

func (s *Stack[T]) Print() {
	fmt.Println("Stack State")
	// print in reverse since the stack works that way
	i2 := s.Len() - 1
	for i := i2; i >= 0; i-- {
		fmt.Println("=>", i2-i, "=", s.arr[i])
	}
}
func (s *Stack[T]) MustPop() T {
	val, err := s.Pop()
	if err != nil {
		log.Fatal(red.Sprintf("call stack is empty"))
	}
	return val
}

var ErrStackEmpty = errors.New("stack is empty")

func (s *Stack[T]) Pop() (T, error) {
	if len(s.arr) == 0 {
		var zero T
		return zero, ErrStackEmpty
	}

	i := len(s.arr) - 1
	val := s.arr[i]
	s.arr = s.arr[:i]
	return val, nil
}

func (s *Stack[T]) LoadAll(arr []T) {
	s.arr = arr
}
