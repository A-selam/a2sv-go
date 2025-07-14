package main

import (
	"fmt"
	"strings"
	"unicode"
)

func cleanUp(str string) string{
	new_str := ""
	for _, r := range str{
		if unicode.IsPunct(r) || unicode.IsSymbol(r){
			continue
		}
		new_str += string(r) 
	}
	return strings.ToLower(new_str)
}

// wordFrequencyCount returns a map of each word's frequency in the input string.
// The function is case-insensitive and ignores punctuation and symbols.
func wordFrequencyCount(str string) map[string]int {
	cleaned := cleanUp(str)
	strings_slice := strings.Fields(cleaned)
	count := make(map[string]int)

	for _, word := range strings_slice{
		count[word]++
	}

	return count
}

func main() {
	text := "Go Go GO! Is Go fun? Yes, go is fun."
	result := wordFrequencyCount(text)
	fmt.Println(result)
}
