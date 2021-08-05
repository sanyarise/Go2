package main

import (
	"fmt"
	"reflect"
	"strings"
)

func inQuery(q string, args ...interface{}) (string, []interface{}) {

	var slic []interface{}
	inString := "("

	for g := 1; g <= reflect.ValueOf(args).Len(); g++ {
		if g == reflect.ValueOf(args).Len() {
			inString = inString + "?)"
			break
		}
		inString = inString + "?,"
	}
	updateQ := strings.Replace(q, "(?)", inString, -1)

	refArgs := reflect.ValueOf(args)
	if refArgs.Kind() == reflect.Slice {
		sliceLen := refArgs.Len()
		for i := 0; i < sliceLen; i++ {
			m := refArgs.Index(i).Interface()
			v := m.(interface{})
			t := reflect.TypeOf(m)
			if t.Kind() == reflect.Slice {
				a := reflect.ValueOf(m)
				for k := 0; k < a.Len(); k++ {
					innerSlice := a.Index(k).Interface().(interface{})
					slic = append(slic, innerSlice)
				}

			} else {
				slic = append(slic, v)
			}
		}
	}
	return updateQ, slic
}

func main() {
	query, args := inQuery("SELECT * FROM table WHERE deleted = ? AND id IN(?) AND count < ?", false, []int{1, 6, 234}, 555)
	fmt.Println(query, args)
}
