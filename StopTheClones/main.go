package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"StopTheClones/clones"
	"StopTheClones/files"
)

var fs files.FileIO

type Filehash struct {
	hash string
	path string
	size int64
}

// FindFilesInPath рекурсивно обходит дерево папок и
// создает хэш для каждой из них
func FindFilesInPath(rootDir string) <-chan Filehash {
	fileChannel := make(chan Filehash)

	go func() {
		defer close(fileChannel)
		fs.Walk(rootDir, func(path string, info files.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			fileChannel <- Filehash{"", path, info.Size()}
			return nil
		})
	}()

	return fileChannel
}

const fileHashSize = 4096

func emitClones(fileChannel <-chan Filehash, clonesChannel chan<- clones.Clone) {
	hashMap := make(map[int64][]Filehash)
	for hf := range fileChannel {
		if fileslice, ok := hashMap[hf.size]; ok {
			hf.hash, _ = fs.MD5HashFile(hf.path, fileHashSize)
			for i := range fileslice {
				f := &fileslice[i]
				if f.hash == "" {
					f.hash, _ = fs.MD5HashFile(f.path, fileHashSize)
				}
				if hf.hash == f.hash && fs.FilesBytesAreEqual(hf.path, f.path) {
					clonesChannel <- clones.Clone{Value1: f.path, Value2: hf.path}
					break
				}
			}
		}
		hashMap[hf.size] = append(hashMap[hf.size], hf)
	}
}

// FindClones возвращает новый канал, содержащий все дублирующиеся
// файлы, обнаруженные в fileChannel
func FindClones(fileChannel <-chan Filehash) <-chan clones.Clone {
	clonesChannel := make(chan clones.Clone)

	go func() {
		defer close(clonesChannel)
		emitClones(fileChannel, clonesChannel)
	}()

	return clonesChannel
}

// GetCloneFileDeleter возвращает функцию, которая удаляет дублирующиеся файлы с диска
// и выводит сообщение об этом в консоль
func GetCloneFileDeleter(writer io.Writer) clones.CloneHandler {
	return func(d clones.Clone) {
		path := d.Value2
		fmt.Fprintf(writer, "DELETING: %s\n", path)
		fs.Delete(d.Value2)
	}
}

// ProcessCloneFiles проверяет все файлы, обнаруженные рекурсивно в папке, на совпадение хэша, используя
// соответствующую функцию CloneHandler
func ProcessCloneFiles(dir string, cloneHandler clones.CloneHandler) {
	clones.ApplyFuncToChan(
		FindClones(
			FindFilesInPath(dir)), cloneHandler)
}

func SetFS(newFS files.FileIO) {
	fs = newFS
}

// Данная программа получает на вход через консоль путь к папке, в которой
// производит поиск всех повторяющихся файлов и выводит их в консоль.
// При запуске с флагом -c или без флагов выводит список дубликатов в консоль.
// При запуске с флагом -d по окончании поиска выводится запрос на подтверждение
// удаления дубликатов. При положительном ответе производится удаление дублирующихся
// файлов.
func main() {
	csv := flag.Bool("c", false, "Print duplicate values as a CSV to the console")
	del := flag.Bool("d", false, "Delete all duplicate values")

	flag.Parse()

	cloneHandler := clones.GetWriter(os.Stdout)
	if *csv {
		cloneHandler = clones.GetCSVWriter(os.Stdout)
	} else if *del {
		for {
			fmt.Println("Are you sure? Files will be deleted without possible to recover!(Y/N)")
			var isDel string
			fmt.Scan(&isDel)
			if isDel == "N" || isDel == "n" {
				cloneHandler = clones.GetCSVWriter(os.Stdout)
				break
			} else if isDel == "Y" || isDel == "y" {
				cloneHandler = GetCloneFileDeleter(os.Stdout)
				break
			}
			continue
		}

	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [options] directory\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Calling without any options does a dry run and lists the files to be deleted")
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}

	dir := flag.Arg(0)

	SetFS(files.FS{})
	ProcessCloneFiles(dir, cloneHandler)
}

func ExampleMain() {

	csv := flag.Bool("c", false, "Print duplicate values as a CSV to the console")
	del := flag.Bool("d", false, "Delete all duplicate values")

	flag.Parse()

	cloneHandler := clones.GetWriter(os.Stdout)
	if *csv {
		cloneHandler = clones.GetCSVWriter(os.Stdout)
	} else if *del {
		for {
			fmt.Println("Are you sure? Files will be deleted without possible to recover!(Y/N)")
			var isDel string
			fmt.Scan(&isDel)
			if isDel == "N" || isDel == "n" {
				cloneHandler = clones.GetCSVWriter(os.Stdout)
				break
			} else if isDel == "Y" || isDel == "y" {
				cloneHandler = GetCloneFileDeleter(os.Stdout)
				break
			}
			continue
		}

	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [options] directory\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Calling without any options does a dry run and lists the files to be deleted")
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}

	dir := flag.Arg(0)

	SetFS(files.FS{})
	ProcessCloneFiles(dir, cloneHandler)
}
