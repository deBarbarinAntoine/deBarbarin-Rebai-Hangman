package ProjetHangman

import (
	"fmt"
	"math/rand"
	"time"
)

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

func displayWord(word []rune, wordDisplay []rune, char rune) []rune {
	for i, r := range word {
		if r == char {
			wordDisplay[i*2] = char - 32
			nbLettersFound++
		}
	}
	return wordDisplay
}
/*  func Randplay() {

	rand.Seed(time.Now().Unix())

	//Generate a random character between lowercase a to z
	randomChar := 'a' + rune(rand.Intn(26))
	fmt.Println(string(randomChar))

	//Generate a random character between uppercase A and Z
	randomChar = 'A' + rune(rand.Intn(26))
	fmt.Println(string(randomChar))

	//Generate a random character between uppercase A and Z  and lowercase a to z
	randomInt := rand.Intn(2)
	if randomInt == 1 {
		randomChar = 'A' + rune(rand.Intn(26))
	} else {
		randomChar = 'a' + rune(rand.Intn(26))
	}
	fmt.Println(string(randomChar)) 
} */ 
func Randplay(n int) string {
    var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
