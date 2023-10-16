package main

import "fmt"


func CheckWord(word string) {
	var input string
	fmt.Print("Rentre un mot : ")
	fmt.Scanln(&input)
	if input == word {
		fmt.Print("Mot identique.")
	} else {
		fmt.Print("Mot non identique.")
	}
}

func displayWord (word []rune, wordDisplay []rune, char rune) []rune { 
	for i, r := range word {
		if r == char {
			wordDisplay[i] = char
		}
	}
	return wordDisplay
} 