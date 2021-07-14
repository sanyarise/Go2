package main

import (
	"fmt"
	"os"
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println("recovered", v)
		}
	} ()
	file, err := os.Create("New_file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, _ = fmt.Fprintln(file, "Text for new file")
}