package ProjetHangman

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
)

type Game struct {
	Name        string
	Score       int
	Word        string
	Difficulty  int
	Dictionnary string
}

var savedGames []Game

func checkWord(word []rune, try string) bool {
	return try == string(word)
}

func displayWord(word []rune, wordDisplay []rune, char rune) []rune {
	for i, r := range word {
		if r == char {
			wordDisplay[i*2] = char - 32
			nbLettersFound++
			score += 10
		}
	}
	return wordDisplay
}

func revealWord(word, wordDisplay []rune) []rune {
	for i, r := range word {
		wordDisplay[i*2] = r - 32
	}
	return wordDisplay
}

func nbRemainingLetters(wordDisplay []rune) int {
	var result int
	for _, char := range wordDisplay {
		if char == '_' {
			result++
		}
	}
	return result
}

/*
	  func Randplay() {

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
	}
*/
func Randplay(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func saveGame() {
	currentGame := Game{
		Name:        name,
		Score:       score,
		Word:        string(word),
		Difficulty:  difficulty,
		Dictionnary: "Scrabble",
	}
	newEntry, err := json.Marshal(currentGame)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile("../Files/scores.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	newEntry = append([]byte{',', '\n'}, newEntry...)
	_, err = file.Write(newEntry)
	if err != nil {
		log.Fatal(err)
	}
}

func retreiveSavedGames() {
	savedEntries, err := os.ReadFile("../Files/scores.txt")
	if err != nil {
		fmt.Println("Aucune sauvegarde détectée...")
		return
	}
	savedEntries = append([]byte{'[', '\n'}, savedEntries...)
	savedEntries = append(savedEntries, '\n', ']')
	err = json.Unmarshal(savedEntries, &savedGames)
	if err != nil {
		fmt.Println("Erreur de récupération des données...")
		fmt.Println()
		fmt.Println("Données récupérées :")
		fmt.Println(string(savedEntries))
		log.Fatal(err)
	}
}

func filterTopTenGames() []Game {
	retreiveSavedGames()
	sort.SliceStable(savedGames, func(i, j int) bool { return savedGames[i].Score > savedGames[j].Score })
	return savedGames
}

func toStringDifficulty(difficulty int) string {
	switch difficulty {
	case EASY:
		return "Facile"
	case MEDIUM:
		return "Intermédiaire"
	case DIFFICULT:
		return "Difficile"
	case LEGENDARY:
		return "Légendaire"
	default:
		return "Inconnu"
	}
}
