package snapshot

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"

	"github.com/itsliamegan/watch/changeset"
)

type Snapshot struct {
	dir        string
	signatures map[string][]byte
}

func Take(dir string) (*Snapshot, error) {
	signatures, err := computeDirSignatures(dir)
	if err != nil {
		return nil, err
	}

	return &Snapshot{dir, signatures}, nil
}

func Compare(old, new *Snapshot) *changeset.ChangeSet {
	changes := changeset.New()

	for path, oldSignature := range old.signatures {
		if newSignature, found := new.signatures[path]; found {
			if !bytes.Equal(oldSignature, newSignature) {
				changes.Set(path, changeset.OperationUpdate)
			}
		} else {
			changes.Set(path, changeset.OperationRemove)
		}
	}

	for path, _ := range new.signatures {
		if _, found := old.signatures[path]; !found {
			changes.Set(path, changeset.OperationCreate)
		}
	}

	return changes
}

func computeDirSignatures(dir string) (map[string][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	signatures := make(map[string][]byte)

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			innerSignatures, err := computeDirSignatures(path)
			if err != nil {
				return nil, err
			}

			for innerPath, signature := range innerSignatures {
				signatures[innerPath] = signature
			}
		} else {
			signature, err := computeFileSignature(path)
			if err != nil {
				return nil, err
			}

			signatures[path] = signature
		}
	}

	return signatures, nil
}

func computeFileSignature(file string) ([]byte, error) {
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	digest := sha256.New()
	_, err = io.Copy(digest, fp)
	if err != nil {
		return nil, err
	}

	signature := make([]byte, 0, digest.Size())
	signature = digest.Sum(signature)

	return signature, nil
}
