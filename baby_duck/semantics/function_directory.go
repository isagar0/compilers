package semantics

import (
	"fmt"
	"sort"
)

type Function struct {
	Name       string
	ReturnType string
	Parameters []string
	Vars       *VariableTable
}

type FunctionDirectory struct {
	directory map[string]*Function
}

func NewFunctionDirectory() *FunctionDirectory {
	return &FunctionDirectory{
		directory: make(map[string]*Function),
	}
}

func (fd *FunctionDirectory) Put(name, returnType string, params []string) error {
	if _, exists := fd.directory[name]; exists {
		return fmt.Errorf("Function '%s' already declared", name)
	}
	fd.directory[name] = &Function{
		Name:       name,
		ReturnType: returnType,
		Parameters: params,
		Vars:       NewVariableTable(),
	}
	return nil
}

func (fd *FunctionDirectory) Get(name string) (*Function, bool) {
	f, ok := fd.directory[name]
	return f, ok
}

func (fd *FunctionDirectory) Remove(name string) {
	delete(fd.directory, name)
}

func (fd *FunctionDirectory) IsEmpty() bool {
	return len(fd.directory) == 0
}

func (fd *FunctionDirectory) Size() int {
	return len(fd.directory)
}

func (fd *FunctionDirectory) Keys() []string {
	keys := make([]string, 0, len(fd.directory))
	for k := range fd.directory {
		keys = append(keys, k)
	}
	return keys
}

func (fd *FunctionDirectory) PrintOrdered() {
	keys := fd.Keys()
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("Function %s -> return: %s, params: %v\n", k, fd.directory[k].ReturnType, fd.directory[k].Parameters)
	}
}
