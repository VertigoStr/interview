## Интересные задачи на go

#### 1. Slice, capasity, append.
Что выведет следующий код?
```go
package main

import "fmt"

func main() {
    slice := []int{1, 2}
    slice = append(slice, 3)
    x := append(slice, 4)
    x = append(slice, 5)
    y := append(slice, 6)
    fmt.Println(x, y)
}
```
<details>
    <summary>Ответ</summary>
    [1 2 3 6] [1 2 3 6]
</details>


#### 2. Stack
```go
package main

import (
 "fmt"
 "sync"
)

// Структура для представления элемента стека
type StackElement struct {
 Value int
}

// Структура для представления стека
type Stack struct {
 elements []StackElement
 mutex    sync.Mutex
}

// Метод для добавления элемента в стек
func (s *Stack) Push(value int) {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 s.elements = append(s.elements, StackElement{Value: value})
}

// Метод для удаления элемента из стека
func (s *Stack) Pop() (int, bool) {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 if len(s.elements) == 0 {
  return 0, false
 }
 last := s.elements[len(s.elements)-1]
 s.elements = s.elements[:len(s.elements)-1]
 return last.Value, true
}

// Метод для проверки пуста ли очередь
func (s *Stack) IsEmpty() bool {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 return len(s.elements) == 0
}

// Метод для поиска максимального элемента в стеке
func (s *Stack) FindMax() int {
 s.mutex.Lock()
 defer s.mutex.Unlock()
 max := s.elements[len(s.elements)-1].Value
 for _, v := range s.elements[:len(s.elements)-1] {
  if v.Value > max {
   max = v.Value
  }
 }
 return max
}

func main() {
 var wg sync.WaitGroup
 wg.Add(1)
 go func() {
  defer wg.Done()
  stack := &Stack{}
  for i := 0; i < 10; i++ {
   stack.Push(i)
  }
  fmt.Println("Максимальный элемент в стеке:", stack.FindMax())
 }()
 wg.Wait()
}
```
