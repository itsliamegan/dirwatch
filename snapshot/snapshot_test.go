package snapshot

import (
	"os"
	"testing"

	"github.com/itsliamegan/watch/changeset"

	"github.com/itsliamegan/assert"
)

func TestComparesCreatedFile(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.MkdirTemp("", "dir")
	defer os.RemoveAll(dir)

	before, _ := Take(dir)

	file, _ := os.CreateTemp(dir, "file")
	defer os.Remove(file.Name())
	defer file.Close()

	after, _ := Take(dir)
	changes := Compare(before, after)

	assert.Equal(changes.Get(file.Name()), changeset.OperationCreate)
}

func TestComparesUpdatedFile(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.MkdirTemp("", "dir")
	defer os.RemoveAll(dir)
	file, _ := os.CreateTemp(dir, "file")
	defer os.Remove(file.Name())
	defer file.Close()

	file.WriteString("before")
	before, _ := Take(dir)
	file.WriteString("after")
	after, _ := Take(dir)
	changes := Compare(before, after)

	assert.Equal(changes.Get(file.Name()), changeset.OperationUpdate)
}

func TestComparesRemovedFile(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.MkdirTemp("", "dir")
	defer os.RemoveAll(dir)
	file, _ := os.CreateTemp(dir, "file")
	defer file.Close()

	before, _ := Take(dir)
	os.Remove(file.Name())
	after, _ := Take(dir)
	changes := Compare(before, after)

	assert.Equal(changes.Get(file.Name()), changeset.OperationRemove)
}

func TestComparesMultipleChanges(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.MkdirTemp("", "dir")
	defer os.RemoveAll(dir)
	file1, _ := os.CreateTemp(dir, "file1")
	defer file1.Close()

	before, _ := Take(dir)

	os.Remove(file1.Name())
	file2, _ := os.CreateTemp(dir, "file2")
	defer os.Remove(file2.Name())
	defer file2.Close()

	after, _ := Take(dir)
	changes := Compare(before, after)

	assert.Equal(changes.Get(file1.Name()), changeset.OperationRemove)
	assert.Equal(changes.Get(file2.Name()), changeset.OperationCreate)
}

func TestComparesNestedChanges(t *testing.T) {
	assert := assert.New(t)

	dir, _ := os.MkdirTemp("", "dir")
	defer os.RemoveAll(dir)
	innerDir, _ := os.MkdirTemp(dir, "dir")
	defer os.RemoveAll(innerDir)

	before, _ := Take(dir)

	file, _ := os.CreateTemp(innerDir, "file")
	defer os.Remove(file.Name())
	defer file.Close()

	after, _ := Take(dir)
	changes := Compare(before, after)

	assert.Equal(changes.Get(file.Name()), changeset.OperationCreate)
}
