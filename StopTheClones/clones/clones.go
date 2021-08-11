package clones

import (
	"fmt"
	"io"
)

// Clone используется для логического сравнения
// двух путей
type Clone struct {
	Value1 string
	Value2 string
}

// CloneHandler функция, производящая обработку дубликатов
type CloneHandler func(Clone)

// GetWriter возвращает функцию для печати содержимого дубликатов в writer
func GetWriter(writer io.Writer) CloneHandler {
	return func(d Clone) {
		fmt.Fprintln(writer, d.Value2)
	}
}

// GetCSVWriter возвращает функцию, которая печатает повторяющееся значение и значение, которое оно дублирует, в формате csv
func GetCSVWriter(writer io.Writer) CloneHandler {
	return func(d Clone) {
		fmt.Fprintf(writer, "\"%s\",\"%s\"\n", d.Value1, d.Value2)
	}
}

// ApplyFuncToChan итерируется по каналу с повторяющимися файлами и применяет обработчик для каждого из них
func ApplyFuncToChan(clones <-chan Clone, handler CloneHandler) {
	for d := range clones {
		handler(d)
	}
}
