package main

import (
	"fmt"
	"strings"
	"unicode"
)

func cleanUp2(str string) string{
	new_str := ""
	for _, r := range str{
		if unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r){
			continue
		}
		new_str += string(r) 
	}
	return strings.ToLower(new_str)
}

func isPalindrome(str string)bool{
	cleaned := cleanUp2(str)
	left := 0
	right := len(cleaned)-1
	for left < right{
		if cleaned[left] != cleaned[right]{
			return false
		}
		left++
		right--
	}
	return true
}

func main() {
	text := "ab Ba"
	result := isPalindrome(text)
	fmt.Println(result)
}
