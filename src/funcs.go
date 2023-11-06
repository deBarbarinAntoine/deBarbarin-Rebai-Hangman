package ProjetHangman

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Game struct {
	Name        string
	Score       int
	Word        string
	Difficulty  int
	Dictionnary string
}

type Parameters struct {
	Name            string
	DictionaryPath  string
	ColorBorder     Color
	ColorTitle      Color
	ColorOptions    Color
	ColorPointingAt Color
}

var savedGames []Game

// checkWord checks if the player found the word. Returns true if he found it and false if not.
func checkWord(word []rune, try string) bool {
	return try == string(word)
}

// Function that changes the wordDisplay to replace the '_' character with the rune played if it is in the word.
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

// revealWord reveals all runes in wordDisplay.
func revealWord(word, wordDisplay []rune) []rune {
	for i, r := range word {
		wordDisplay[i*2] = r - 32
	}
	return wordDisplay
}

// nbRemainingLetters returns the number of letters still not found in the word.
func nbRemainingLetters(wordDisplay []rune) int {
	var result int
	for _, char := range wordDisplay {
		if char == '_' {
			result++
		}
	}
	return result
}

// retreiveWords retreive the words from the selected dictionary.
func retreiveWords() {
	if dictionaryPath == "" {
		dictionaryPath = "../Files/Dictionaries/ods5.txt"
	}
	content, err := os.ReadFile(dictionaryPath)
	if err != nil {
		log.Fatal(err)
	}
	if checkDictionary() {
		words = strings.Split(string(content), "\n")
	} else {
		fmt.Println(colorCode(Red), "Erreur d'acquisition des mots du dictionnaire", CLEARCOLOR)
		time.Sleep(time.Second * 2)
	}
}

// retreiveHangman retreives the hangman in /Files/hangman.txt and stores it in hangman.
func retreiveHangman() {
	hangman = append(hangman[0:0])
	content, err := os.ReadFile("../Files/hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	var line int
	var str string
	for _, char := range content {
		str += string(char)
		if char == '\n' {
			line++
		}
		if line == 8 {
			hangman = append(hangman, str)
			str = ""
			line = 0
		}
	}
}

// checkRune checks if the rune played is already played, correct or incorrect.
func checkRune(char rune) int {
	for _, r := range runesPlayed {
		if r == char {
			return ALREADYPLAYED
		}
	}
	for _, r := range word {
		if r == char {
			return CORRECTRUNE
		}
	}
	return INCORRECTRUNE
}

// clearGameData clears the previous' game's data to start a new one.
func clearGameData() {
	nbErrors = 0
	nbLettersFound = 0
	score = 0
	runesPlayed = append(runesPlayed[0:0])
}

// hint reveal a random rune in wordDisplay.
func hint(wordDisplay []rune) []rune {
	if difficulty != LEGENDARY {
		i := rand.Intn(len(word) - 1)
		char := word[i]
		wordDisplay[i*2] = char - 32
	}
	return wordDisplay
}

// chooseWord chooses randomly a word from words (the dictionary's words' list) according to the difficulty set previously.
func chooseWord(difficulty int) []rune {
	var possibleWords []string
	for _, str := range words {
		if len(str) >= difficulty-2 && len(str) <= difficulty {
			possibleWords = append(possibleWords, str)
		}
		if difficulty == LEGENDARY {
			if len(str) > difficulty {
				possibleWords = append(possibleWords, str)
			}
		}
	}
	if len(possibleWords) < 10 {
		var i int
		for _, str := range words {
			i++
			if len(str) == difficulty-i-2 || len(str) == difficulty+i {
				possibleWords = append(possibleWords, str)
				if len(possibleWords) > 10 {
					break
				}
			}
		}
	}
	return []rune((possibleWords[rand.Intn(len(possibleWords)-1)]))
}

// dictionaryName returns the name of the dictionary.
func dictionaryName() string {
	switch dictionaryPath {
	case "../Files/Dictionaries/ods5.txt":
		return "Scrabble"
	case "../Files/Dictionaries/ospd3_expurgated.txt":
		return "Anglais"
	case "../Files/Dictionaries/italiano.txt":
		return "Italien"
	default:
		return "Personnalisé"
	}
}

// saveGame saves the current game in /Files/scores.txt.
func saveGame() {
	currentGame := Game{
		Name:        name,
		Score:       score,
		Word:        string(word),
		Difficulty:  difficulty,
		Dictionnary: dictionaryName(),
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

// retreiveSavedGames retreive all saved games present in /Files/scores.txt and put it in savedEntries.
func retreiveSavedGames() {
	savedEntries, err := os.ReadFile("../Files/scores.txt")
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucune sauvegarde détectée...", CLEARCOLOR)
		return
	}
	savedEntries = append([]byte{'[', '\n'}, savedEntries...)
	savedEntries = append(savedEntries, '\n', ']')
	err = json.Unmarshal(savedEntries, &savedGames)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orangered), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	}
}

// sortTopTenGames sort the saved games by score in decreasing order.
func sortTopTenGames() []Game {
	retreiveSavedGames()
	sort.SliceStable(savedGames, func(i, j int) bool { return savedGames[i].Score > savedGames[j].Score })
	return savedGames
}

// toStringDifficulty returns the name of the difficulty.
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

// checkDictionary checks if the dictionary is usable or not and changes the case. It tries to remove the accents, but the function doesn't really work unfortunately.
func checkDictionary() bool {
	for i, str := range words {
		words[i] = removeAccents(words[i])
		words[i] = strings.ToLower(str)
		str = words[i]
		for _, char := range str {
			if char < 'a' || char > 'z' {
				return false
			}
		}
	}
	return true
}

// removeAccents removes the accents. Unfortunately, it doesn't really work.
func removeAccents(str string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(transformer, str)
	if err != nil {
		log.Fatal(err)
	}
	for i, char := range result {
		if char == 'ñ' {
			result = string([]rune(result)[:i]) + "n" + string([]rune(result)[i+1:])
		}
	}
	return result
}

// chargeParameters retreive the parameters present in /Files/config.txt and changes all corresponding variables.
func chargeParameters() {
	var savedParameters Parameters
	savedEntries, err := os.ReadFile("../Files/config.txt")
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucune fichier de configuration détecté...", CLEARCOLOR)
		time.Sleep(time.Second * 1)
		return
	}
	err = json.Unmarshal(savedEntries, &savedParameters)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "Il est conseillé de supprimer le fichier config.txt\n    dans le dossier Files afin de résoudre le problème.", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orange), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	} else {
		name = savedParameters.Name
		dictionaryPath = savedParameters.DictionaryPath
		colorBorder = savedParameters.ColorBorder
		colorTitle = savedParameters.ColorTitle
		colorOptions = savedParameters.ColorOptions
		colorPointingAt = savedParameters.ColorPointingAt
	}
}

// saveParameters saves all current parameters in /Files/config.txt for later use.
func saveParameters() {
	currentParameters := Parameters{
		Name:            name,
		DictionaryPath:  dictionaryPath,
		ColorBorder:     colorBorder,
		ColorTitle:      colorTitle,
		ColorOptions:    colorOptions,
		ColorPointingAt: colorPointingAt,
	}
	newEntry, err := json.Marshal(currentParameters)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("../Files/config.txt", newEntry, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
