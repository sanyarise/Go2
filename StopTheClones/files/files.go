package files

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"StopTheClones/bytesequal"
)

type FileInfo interface {
	IsDir() bool
	Size() int64
}

// Функция, вызываемая для каждого элемента, найденного Walk
type WalkFn = func(path string, info FileInfo, err error) error

// Интерфейс для действий с файлами
type FileIO interface {
	Delete(string) error
	FilesBytesAreEqual(string, string) bool
	MD5HashFile(string, int64) (string, error)
	Walk(string, WalkFn) error
}


type FS struct{}

// readfile читает файл по переданному path
func readfile(path string) []byte {
	if data, err := ioutil.ReadFile(path); err == nil {
		return data
	}
	return nil
}

//производится удаление папки.
func (FS) Delete(path string) error {
	return os.Remove(path)
}

// FilesBytesAreEqual возвращает true если два файла полностью совпадают
func (FS) FilesBytesAreEqual(path1 string, path2 string) bool {
	b1 := readfile(path1)
	b2 := readfile(path2)
	return bytesequal.BytesAreEqual(b1, b2)
}

// MD5HashFile создает хэш для файла
func (FS) MD5HashFile(path string, hashSize int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return bytesequal.MD5Hash(file, hashSize)
}

// Обходит все пути в исходной папке и вызывает WalkFn для каждого найденного элемента
func (FS) Walk(root string, walkFn WalkFn) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return walkFn(path, info, err)
	})
}
