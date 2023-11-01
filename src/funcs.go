package ProjetHangman

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
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

func sortTopTenGames() []Game {
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
