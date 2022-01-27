package changeset

import (
	"fmt"
)

type ChangeSet struct {
	operations map[string]Operation
}

type Change struct {
	Path      string
	Operation Operation
}

type Operation uint8

const (
	OperationNone Operation = iota
	OperationCreate
	OperationUpdate
	OperationRemove
)

func New() *ChangeSet {
	return &ChangeSet{operations: make(map[string]Operation)}
}

func (changes *ChangeSet) All() []Change {
	all := make([]Change, 0, len(changes.operations))

	for path, operation := range changes.operations {
		all = append(all, Change{path, operation})
	}

	return all
}

func (changes *ChangeSet) Get(path string) Operation {
	return changes.operations[path]
}

func (changes *ChangeSet) Set(path string, operation Operation) {
	changes.operations[path] = operation
}

func (change Change) String() string {
	return fmt.Sprintf("%s\t%s", change.Operation, change.Path)
}

func (operation Operation) String() string {
	switch operation {
	case OperationCreate:
		return "create"
	case OperationUpdate:
		return "update"
	case OperationRemove:
		return "remove"
	default:
		panic("invalid operation")
	}
}
