package main

import (
	"bytes"
	"fmt"
)

func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func commaNoRecursion(s string) string {
	var buf bytes.Buffer
	offset := 0
	for i := len(s) - 1; i >= 0; i-- {
		if 0 < offset && 0 == offset%3 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
		offset += 1
	}
	return Reverse(buf.String())
}

func main() {
	fmt.Printf("123456: %s\n", comma("123456"))
	fmt.Printf("123456: %s\n", commaNoRecursion("123456"))
}
