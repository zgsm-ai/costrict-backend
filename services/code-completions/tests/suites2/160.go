package main

import (
    "container/list"
    "fmt"
)

type Stack struct {
    items *list.List
}

func NewStack() *Stack {
    return &Stack{
        items: list.New(),
    }
}

func (s *Stack) Push(value interface{}) {
    s.items.PushBack(value)
}

func (s *Stack) Pop() interface{} {
    if s.items.Len() == 0 {
        return nil
    }
    
    element := s.items.Back()
    <｜fim▁hole｜>
    
    return element.Value
}

func (s *Stack) Peek() interface{} {
    if s.items.Len() == 0 {
        return nil
    }
    return s.items.Back().Value
}

func (s *Stack) IsEmpty() bool {
    return s.items.Len() == 0
}

func (s *Stack) Size() int {
    return s.items.Len()
}

func main() {
    stack := NewStack()
    
    stack.Push(10)
    stack.Push(20)
    stack.Push(30)
    
    fmt.Println("Stack size:", stack.Size())
    fmt.Println("Top element:", stack.Peek())
    
    popped := stack.Pop()
    fmt.Println("Popped element:", popped)
    fmt.Println("Stack size after pop:", stack.Size())
}