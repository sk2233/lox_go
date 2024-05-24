/*
@author: sk
@date: 2024/3/16
*/
package main

import "fmt"

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Report(line int, where string, msg string) {
	fmt.Printf("[line %d] Error %s : %s\n", line, where, msg)
}

func Error(line int, msg string) {
	Report(line, "", msg)
}

func IsDigit(ch uint8) bool {
	return ch >= '0' && ch <= '9'
}

func IsAlpha(ch uint8) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	return ch == '_'
}
