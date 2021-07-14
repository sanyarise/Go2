package main

import (
	"fmt"
	"time"
	"errors"	
)
type MyError struct {
	time string
	text string
}
func New(text string) error {
	return &MyError{
		time: timeStamp(),
		text: text,
	}
}

func (m *MyError) Error() string {
	return fmt.Sprintf("Time: %s\nText: %s", m.time, m.text)
}

func timeStamp()string {
	time := time.Now()
	return time.Format("15:04:05")
}

func main() {
	defer func() {
		if v:=recover(); v!=nil {
			fmt.Println("Внимание! Произошла паника!", v)
		}
	}()
	var err error
	var s int
	var t int
	var v int
	err = errors.New("Ошибка!")
	fmt.Println(err)
	err = New("Ошибка вызвана по ошибке!")
	fmt.Println(err)
	s = 60
	v = s/t
	fmt.Printf("Скорость равна %d", v)
}
