package main

import (
	"fmt"
	"sort"
	"strings"
)

func sortChars(s string) string {
	tmp := strings.Split(s, "")
	sort.Strings(tmp)
	return strings.Join(tmp, "")
}

func isAnagram(lhs, rhs string) bool {
	return sortChars(lhs) == sortChars(rhs)
}

func main() {
	fmt.Printf("%s, %s, %t\n", "abc", "bca", isAnagram("abc", "bca"))
	fmt.Printf("%s, %s, %t\n", "abc", "bcad", isAnagram("abc", "bcad"))
}
