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
	if err != nil {
		fmt.Println("Test Case 2 Error:", err)
	} else {
		fmt.Println("Test Case 2:", stack.items, "Elemento eliminado: ", item) // Esperado: [10 20] Elemento eliminado: 30
	}

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
	if err != nil {
		fmt.Println("Test Case 2 Error:", err)
	} else {
		fmt.Println("Test Case 2:", queue.items, "Elemento eliminado: ", item) // Esperado: [20 30] Elemento eliminado:  10
	}

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

func DictionaryTests() {
	// Test Case 1: Put
	dict := NewDictionary()
	dict.Put("a", 10)
	dict.Put("b", 20)
	dict.Put("c", 30)
	fmt.Println("Test Case 1:", dict.items) // Esperado: map[a:10 b:20 c:30]

	// Test Case 2: Get
	value, exists := dict.Get("b")
	fmt.Println("Test Case 2:", value, exists) // Esperado: 20 true

	// Test Case 3: Get
	value, exists = dict.Get("z")
	fmt.Println("Test Case 3:", value, exists) // Esperado: 0 false

	// Test Case 4: Remove
	dict.Remove("b")
	fmt.Println("Test Case 4:", dict.items) // Esperado: map[a:10 c:30]

	// Test Case 5: Remove
	dict.Remove("z")
	fmt.Println("Test Case 5:", dict.items) // Esperado: Sin cambios

	// Test Case 6: Keys
	fmt.Println("Test Case 6:", dict.Keys()) // Esperado: [a c]

	// Test Case 7: Size
	size := dict.Size()
	fmt.Println("Test Case 7:", size) // Esperado: 2

	// Test Case 8: IsEmpty
	isEmpty := dict.IsEmpty()
	fmt.Println("Test Case 8:", isEmpty) // Esperado: false

	// Test Case 9: IsEmpty
	dict.Remove("a")
	dict.Remove("c")
	fmt.Println("Test Case 9:", dict.IsEmpty()) // Esperado: true

	// Test Case 10: Put
	dict.Put("x", 50)
	dict.Put("x", 100)
	fmt.Println("Test Case 10:", dict.items["x"]) // Esperado: 100

	// Test Case 11: PrintedOrdered
	dict.Put("y", 200)
	dict.Put("x", 300)
	dict.Put("m", 400)
	dict.Put("n", 500)
	fmt.Println("Test Case 12:")
	dict.PrintOrdered() // Esperado: x: 300 y:200 m: 400 n: 500 en orden de inserción
}

func main() {
	StackTests()
	fmt.Println()
	QueueTests()
	fmt.Println()
	DictionaryTests()
}
