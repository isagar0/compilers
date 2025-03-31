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
	fmt.Println("Test Case 2:", stack.items, err) // Esperado: [10 20]

	// Test case 3: Peek
	item, err = stack.Peek()
	fmt.Println("Test Case 3:", item, err) // Esperado: 20

	// Test case 4: IsEmpty
	isEmpty := stack.IsEmpty()
	fmt.Println("Test Case 4:", isEmpty) // Esperado: false

	// Test case 5: Size
	size := stack.Size()
	fmt.Println("Test Case 5:", size) // Esperado: 2

	// Test case 6: IsEmpty
	_, _ = stack.Pop()
	_, _ = stack.Pop()

	isEmpty = stack.IsEmpty()
	fmt.Println("Test Case 6:", isEmpty) // Esperado: true

	// Test case 7: Pop
	item, err = stack.Pop()
	fmt.Println("Test Case 7:", item, err) // Esperado: "0 stack is empty"

	// Test case 8: Peek
	item, err = stack.Peek()
	fmt.Println("Test Case 8:", item, err) // Esperado: "0 stack is empty"
}

func QueueTests() {
	// Test Case 1: Enqueue
	queue := Queue{}
	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)
	fmt.Println("Test Case 1:", queue.items) // Esperado: [10 20 30]

	// Test Case 2: Dequeue
	item, err := queue.Dequeue()
	fmt.Println("Test Case 2:", queue.items, err) // Esperado: [20 30]

	// Test Case 3: Peek
	item, err = queue.Peek()
	fmt.Println("Test Case 3:", item, err) // Esperado: 20

	// Test Case 4: IsEmpty
	isEmpty := queue.IsEmpty()
	fmt.Println("Test Case 4:", isEmpty) // Esperado: false

	// Test Case 5: Size
	size := queue.Size()
	fmt.Println("Test Case 5:", size) // Esperado: 2

	// Test Case 6: IsEmpty
	_, _ = queue.Dequeue()
	_, _ = queue.Dequeue()
	isEmpty = queue.IsEmpty()
	fmt.Println("Test Case 6:", isEmpty) // Esperado: true

	// Test Case 7: Pop
	item, err = queue.Dequeue()
	fmt.Println("Test Case 7:", item, err) // Esperado: "0 queue is empty"

	// Test Case 8: Peek
	item, err = queue.Peek()
	fmt.Println("Test Case 8:", item, err) // Esperado: "0 queue is empty"

	// Test Case 9: Print
	queue.Enqueue(5)
	queue.Enqueue(15)
	queue.Enqueue(25)
	fmt.Print("Test Case 9: ")
	queue.Print() // Esperado: "5 15 25"
}

func main() {
	StackTests()
	fmt.Println()
	QueueTests()
}
