package main

import "fmt"

type Stack[T any] struct {
	container []T
}

func (this *Stack[T]) Empty() bool {
	return len(this.container) == 0
}

func (this *Stack[T]) Size() int {
	return len(this.container)
}

func (this *Stack[T]) Push(item T) {
	this.container = append(this.container, item)
}

func (this *Stack[T]) Peek() T {
	return this.container[len(this.container) - 1]
}

func (this *Stack[T]) Pop() T {
	length := len(this.container)
	if length == 0 {
		panic(fmt.Errorf("Can't pop from empty stack"))
	}

	lastItem := this.container[length - 1]

	this.container = this.container[:length - 1]

	return lastItem
}
