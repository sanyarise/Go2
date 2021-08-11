package main

import (
	"bytes"
	"fmt"
	"testing"

	"StopTheClones/clones"

	"StopTheClones/files"
)

type mockFS struct {
	delete             func(string) error
	filesBytesAreEqual func(string, string) bool
	md5HashFile        func(string, int64) (string, error)
	walk               func(string, files.WalkFn) error
}

func (fs mockFS) Delete(path string) error {
	return fs.delete(path)
}

func (fs mockFS) FilesBytesAreEqual(path1 string, path2 string) bool {
	return fs.filesBytesAreEqual(path1, path2)
}

func (fs mockFS) MD5HashFile(path string, hashSize int64) (string, error) {
	return fs.md5HashFile(path, hashSize)
}

func (fs mockFS) Walk(root string, walkFn files.WalkFn) error {
	return fs.walk(root, walkFn)
}

type mockFileInfo struct {
	isdir bool
	size  int64
}

func (mock mockFileInfo) IsDir() bool {
	return mock.isdir
}

func (mock mockFileInfo) Size() int64 {
	return mock.size
}

func TestHashFilesInPath(t *testing.T) {
	numberOfFiles := 10
	mockHashValue := "THIS IS A MOCK"
	mockhash := func(path string, hashSize int64) (string, error) {
		return mockHashValue, nil
	}

	t.Run("Hash files", func(t *testing.T) {
		mockWalk := func(root string, walkFn files.WalkFn) error {
			for i := 0; i < numberOfFiles; i++ {
				path := fmt.Sprintf("%s%d", "test", i)
				walkFn(path, mockFileInfo{isdir: false, size: 100}, nil)
			}
			return nil
		}

		SetFS(mockFS{
			md5HashFile: mockhash,
			walk:        mockWalk,
		})

		hashed := FindFilesInPath("Test")
		hashCount := 0
		for h := range hashed {
			hashCount++
			if h.hash != "" {
				t.Errorf("Unexpected hash: %s", h.hash)
			}
		}
		if hashCount != numberOfFiles {
			t.Errorf("expected %d, received %d", numberOfFiles, hashCount)
		}
	})

	t.Run("Skip directories", func(t *testing.T) {
		mockWalk := func(root string, walkFn files.WalkFn) error {
			for i := 0; i < numberOfFiles; i++ {
				path := fmt.Sprintf("%s%d", "test", i)
				walkFn(path, mockFileInfo{isdir: true}, nil)
			}
			return nil
		}

		SetFS(mockFS{
			md5HashFile: mockhash,
			walk:        mockWalk,
		})

		hashed := FindFilesInPath("Test")
		hashCount := 0
		for h := range hashed {
			if h.hash != "" {
				hashCount++
			}
		}
		if hashCount != 0 {
			t.Errorf("expected 0, received %d", hashCount)
		}
	})
}

func TestFindClones(t *testing.T) {
	numberOfClons := 5
	mockFileBytesAreEqual := func(string, string) bool { return true }
	mockHashValue := "THIS IS A MOCK"
	mockhash := func(path string, hashSize int64) (string, error) {
		return mockHashValue, nil
	}

	SetFS(mockFS{
		filesBytesAreEqual: mockFileBytesAreEqual,
		md5HashFile:        mockhash,
	})

	hashChannel := make(chan Filehash)

	go func() {
		for i := 0; i < numberOfClons; i++ {
			hashChannel <- Filehash{"", fmt.Sprintf("%s%d", "test", i), 10}
		}
		close(hashChannel)
	}()

	clonesChannel := FindClones(hashChannel)

	clonCount := 0
	for s := range clonesChannel {
		if s.Value1 != s.Value2 {
			clonCount++
		}
	}
	if clonCount != numberOfClons-1 {
		t.Errorf("expected %d, received %d", numberOfClons-1, clonCount)
	}
}

func TestDeleteClone(t *testing.T) {
	buffer := bytes.Buffer{}

	mockDeletePath := "test"
	mockDelete := func(received string) error {
		if received != mockDeletePath {
			t.Errorf("expected %s, received %s", mockDeletePath, received)
		}
		return nil
	}

	SetFS(mockFS{
		delete: mockDelete,
	})

	delete := GetCloneFileDeleter(&buffer)
	delete(clones.Clone{Value1: "First found", Value2: mockDeletePath})

	expected := "DELETING: test\n"
	received := buffer.String()

	if received != expected {
		t.Errorf("expected %q, received %q", expected, received)
	}
}

func TestProcessCloneFiles(t *testing.T) {
	numberOfFiles := 10
	mockHashValue := "THIS IS A MOCK"
	mockhash := func(path string, hashSize int64) (string, error) {
		return mockHashValue, nil
	}
	mockWalk := func(root string, walkFn files.WalkFn) error {
		for i := 0; i < numberOfFiles; i++ {
			path := fmt.Sprintf("%s%d", "test", i)
			walkFn(path, mockFileInfo{isdir: false}, nil)
		}
		return nil
	}
	mockFileBytesAreEqual := func(p1 string, p2 string) bool { return true }

	SetFS(mockFS{
		filesBytesAreEqual: mockFileBytesAreEqual,
		md5HashFile:        mockhash,
		walk:               mockWalk,
	})

	clonCount := 0
	clonHandler := func(clones.Clone) {
		clonCount++
	}

	ProcessCloneFiles("test", clonHandler)
	if clonCount != numberOfFiles-1 {
		t.Errorf("expected %d, received %d", numberOfFiles-1, clonCount)
	}
}
