package main

import "fmt"

func StackTests() {
	// Test case 1: Push
	stack := Stack{}
	stack.Push(10)
	stack.Push(20)
	stack.Push(30)
	fmt.Println("Test Case 1:", stack.items) // Esperado: [10 20 30]

	// Test case 2: Pop
	item, err := stack.Pop()
	fmt.Println("Test Case 2:", stack.items) // Esperado: [10 20]

	// Test case 3: Peek
	item, err = stack.Peek()
	fmt.Println("Test Case 3:", item) // Esperado: 20

	// Test case 4: IsEmpty()
	isEmpty := stack.IsEmpty()
	fmt.Println("Test Case 4:", isEmpty) // Esperado: false

	// Test case 5: Size()
	size := stack.Size()
	fmt.Println("Test Case 5:", size) // Esperado: 2

	// Test case 6: IsEmpty()
	_, _ = stack.Pop()
	_, _ = stack.Pop()

	isEmpty = stack.IsEmpty()
	fmt.Println("Test Case 6:", isEmpty) // Esperado: true

	// Test case 7: Pop
	item, err = stack.Pop()
	fmt.Println("Test Case 7:", item, err) // Esperado: "0 stack is empty"

	// Test case 8: Peek()
	item, err = stack.Peek()
	fmt.Println("Test Case 8:", item, err) // Esperado: "0 stack is empty"
}

func main() {
	StackTests()
}
