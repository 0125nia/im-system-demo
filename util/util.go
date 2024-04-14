package util

import "fmt"

// ErrPrint judge and print err
func ErrPrint(err error, output string) bool {
	if err != nil {
		fmt.Println(output+" error: ", err)
		return true
	}
	return false
}
