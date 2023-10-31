package ProjetHangman

import (
	"github.com/mattn/go-tty"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var words []string
var word, runesPlayed []rune
var nbLettersFound, nbErrors, score int
var hangman []string
var firstGame = true
var name string
var difficulty int
var toStrDifficulty string

var (
	colorBorder     Color = Teal
	colorTitle      Color = Deepskyblue
	colorOptions    Color = Forestgreen
	colorPointingAt Color = Aquamarine
)

const (
	ALREADYPLAYED = 0
	CORRECTRUNE   = 1
	INCORRECTRUNE = 2
	CORRECTWORD   = 3
	INCORRECTWORD = 4

	CHANGECOLORBORDER           = 3
	CHANGECOLORTITLE            = 4
	CHANGECOLOROPTIONS          = 5
	CHANGECOLOROPTIONPOINTINGAT = 6

	EASY      = 4
	MEDIUM    = 7
	DIFFICULT = 10
	LEGENDARY = 13
)

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func clearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func inputMenu() (x, y int, enter bool) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	var char1ArrowKey, char2ArrowKey bool
	for {
		char, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		switch char {
		case 27: //First rune code for arrows
			char1ArrowKey = true

		case 91: //Second rune code for arrows
			if char1ArrowKey {
				char2ArrowKey = true
			}
		case 68: //Left
			if char1ArrowKey && char2ArrowKey {
				return -1, 0, false
			}

		case 65: //Up
			if char1ArrowKey && char2ArrowKey {
				return 0, 1, false
			}

		case 67: //Right
			if char1ArrowKey && char2ArrowKey {
				return 1, 0, false
			}

		case 66: //Down
			if char1ArrowKey && char2ArrowKey {
				return 0, -1, false
			}

		case 13: //Enter
			return 0, 0, true

		default:
			char1ArrowKey, char2ArrowKey = false, false
		}
	}
}

func createVerticalMenu(cursorAt int, cursor, title string, options ...string) string {
	for {
		clearDisplay3d()
		buildDisplay3d(0, 0, colorTitle, []string{title})
		for i, option := range options {
			if cursorAt == i {
				buildDisplay3d(i+2, 4, colorPointingAt, []string{cursor + "\t" + option})
			} else {
				buildDisplay3d(i+2, 4, colorOptions, []string{"\t" + option})
			}
		}
		showDisplay3d()
		_, y, enter := inputMenu()
		switch {
		case y < 0:
			if cursorAt == len(options)-1 {
				cursorAt = 0
			} else {
				cursorAt++
			}
		case y > 0:
			if cursorAt == 0 {
				cursorAt = len(options) - 1
			} else {
				cursorAt--
			}
		case enter:
			return options[cursorAt]
		}
	}
}

func PrincipalMenu() {
	clearTerminal()
	for {
		switch createVerticalMenu(0, "-->", "------------------------- MENU PRINCIPAL -------------------------", "Nouvelle partie", "Paramètres", "Quitter") {
		case "Nouvelle partie":
			setName()
		case "Paramètres":
			parameters()
		case "Quitter":
			clearTerminal()
			os.Exit(0)
		}
	}
}

func parameters() {
	clearTerminal()
	for {
		switch createVerticalMenu(0, "-->", "--------------------------- PARAMETRES ---------------------------", "Changer la couleur de la bordure", "Changer la couleur des titres", "Changer la couleur des options", "Changer la couleur de l'option sélectionnée", "Retour") {
		case "Changer la couleur de la bordure":
			selectColorFamily(CHANGECOLORBORDER)
		case "Changer la couleur des titres":
			selectColorFamily(CHANGECOLORTITLE)
		case "Changer la couleur des options":
			selectColorFamily(CHANGECOLOROPTIONS)
		case "Changer la couleur de l'option sélectionnée":
			selectColorFamily(CHANGECOLOROPTIONPOINTINGAT)
		case "Retour":
			return
		}
	}
}

func selectColorFamily(option int) {
	clearTerminal()
	var title string
	switch option {
	case CHANGECOLORBORDER:
		title = "---------------------- COULEUR DE LA BORDURE ----------------------"
	case CHANGECOLORTITLE:
		title = "----------------------- COULEUR DES TITRES -----------------------"
	case CHANGECOLOROPTIONS:
		title = "----------------------- COULEUR DES OPTIONS -----------------------"
	case CHANGECOLOROPTIONPOINTINGAT:
		title = "----------------- COULEUR DE LA LIGNE DU CURSEUR -----------------"
	}
	switch createVerticalMenu(0, "-->", title, "Rouge", "Orange", "Jaune", "Vert", "Cyan", "Bleu", "Violet", "Rose", "Blanc", "Gris", "Marron", "Retour") {
	case "Rouge":
		selectColor(Reds, option)
	case "Orange":
		selectColor(Oranges, option)
	case "Jaune":
		selectColor(Yellows, option)
	case "Vert":
		selectColor(Greens, option)
	case "Cyan":
		selectColor(Cyans, option)
	case "Bleu":
		selectColor(Blues, option)
	case "Violet":
		selectColor(Purples, option)
	case "Rose":
		selectColor(Pinks, option)
	case "Blanc":
		selectColor(Whites, option)
	case "Gris":
		selectColor(Grays, option)
	case "Marron":
		selectColor(Browns, option)
	case "Retour":
		break
	}
}

func selectColor(color []Color, option int) {
	clearTerminal()
	var options []string
	var newColor Color
	for _, c := range color {
		options = append(options, colorCode(c)+c.Name)
	}
	options = append(options, "Retour")
	colorName := createVerticalMenu(0, "-->", "--------------------- CHOISISSEZ UNE COULEUR ---------------------", options...)
	if colorName == "Retour" {
		return
	}
	for _, c := range color {
		if colorCode(c)+c.Name == colorName {
			newColor = c
		}
	}
	switch option {
	case CHANGECOLORBORDER:
		colorBorder = newColor
	case CHANGECOLORTITLE:
		colorTitle = newColor
	case CHANGECOLOROPTIONS:
		colorOptions = newColor
	case CHANGECOLOROPTIONPOINTINGAT:
		colorPointingAt = newColor
	}
}

func retreiveWords() {
	content, err := os.ReadFile("../Files/Dictionaries/ods5.txt")
	if err != nil {
		log.Fatal(err)
	}
	words = strings.Split(string(content), "\n")
}

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

func input() (rune, int, bool) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	var char1ArrowKey, char2ArrowKey bool
	for {
		char, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		switch char {
		case 27: //First rune code for arrows
			char1ArrowKey = true

		case 91: //Second rune code for arrows
			if char1ArrowKey {
				char2ArrowKey = true
			}

		case 65, 66, 67, 68: //Arrow
			if char1ArrowKey && char2ArrowKey {
				return 182, CORRECTRUNE, true // returns: ¶, Correct Rune , Arrow
			}

		default:
			char1ArrowKey, char2ArrowKey = false, false
		}
		if char >= 'A' && char <= 'Z' {
			char += 32
		}
		if char >= 'a' && char <= 'z' {
			switch checkRune(char) {
			case CORRECTRUNE:
				runesPlayed = append(runesPlayed, char)
				return char, CORRECTRUNE, false
			case INCORRECTRUNE:
				runesPlayed = append(runesPlayed, char)
				nbErrors++
				score -= 5
				return '\n', INCORRECTRUNE, false
			case ALREADYPLAYED:
				return '\n', ALREADYPLAYED, false
			}
		}
	}
}

func wordInput() (rune, bool, bool) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	var char1ArrowKey, char2ArrowKey bool
	for {
		char, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		switch char {
		case 27: //First rune code for arrows
			char1ArrowKey = true

		case 91: //Second rune code for arrows
			if char1ArrowKey {
				char2ArrowKey = true
			}

		case 65, 66, 67, 68: //Arrow
			if char1ArrowKey && char2ArrowKey {
				return 182, false, true // returns: ¶, No Enter, Arrow
			}

		case 13: //Enter
			return 182, true, false // returns: ¶, Enter, No Arrow

		default:
			char1ArrowKey, char2ArrowKey = false, false
		}
		if char >= 'A' && char <= 'Z' {
			char += 32
		}
		if (char >= 'a' && char <= 'z') || char == 8 {
			return char, false, false
		}
	}
}

func setName() {
	var incorrectName bool
	for {
		buildDisplay3d(0, 0, colorTitle, []string{"----------------------- SAISISSEZ VOTRE NOM -----------------------"})
		buildDisplay3d(2, 2, colorOptions, []string{"Seuls les lettres sont autorisées (sans espace)"})
		if incorrectName {
			buildDisplay3d(4, 2, Red, []string{"Votre nom doit avoir au moins 3 lettres !"})
		}
		buildDisplay3d(6, 4, colorOptions, []string{"Nom :"})
		buildDisplay3d(6, 10, colorPointingAt, []string{name})
		showDisplay3d()
		char, isSet, _ := wordInput()
		if isSet {
			if len(name) > 2 {
				break
			} else {
				incorrectName = true
			}
		}
		if char == 8 {
			if len(name) > 0 {
				name = string([]rune(name)[:len(name)-1])
			}
		} else {
			name += string(char)
		}
	}
	setDifficulty()
	play()
}

func setDifficulty() {
	clearTerminal()
	for {
		switch createVerticalMenu(0, "-->", "---------------------- CHOISIR LA DIFFICULTÉ ----------------------", "Facile", "Intermédiaire", "Difficile", "Légendaire") {
		case "Facile":
			difficulty = EASY
			toStrDifficulty = "Facile"
			return
		case "Intermédiaire":
			difficulty = MEDIUM
			toStrDifficulty = "Intermédiaire"
			return
		case "Difficile":
			difficulty = DIFFICULT
			toStrDifficulty = "Difficile"
			return
		case "Légendaire":
			difficulty = LEGENDARY
			toStrDifficulty = "Légendaire"
			return
		}
	}
}

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

func hint(wordDisplay []rune) []rune {
	if difficulty != LEGENDARY {
		i := rand.Intn(len(word) - 1)
		char := word[i]
		wordDisplay[i*2] = char - 32
	}
	return wordDisplay
}

func play() {
	var hasWon bool
	if firstGame {
		retreiveWords()
		retreiveHangman()
		firstGame = false
	}
	clearTerminal()
	clearGameData()
	word = chooseWord(difficulty)
	wordDisplay := []rune(strings.Repeat("_ ", len(word)))
	hint(wordDisplay)
	var char rune
	status := CORRECTRUNE
	var wordMode bool
	var try string
	for {
		buildDisplay3d(2, 4, colorTitle, []string{"Score : " + colorCode(colorPointingAt) + strconv.Itoa(score)})
		buildDisplay3d(2, 38, colorTitle, []string{"Difficulté : " + colorCode(colorPointingAt) + toStrDifficulty})
		buildDisplay3d(4, 52, colorTitle, strings.Split(hangman[nbErrors], "\n"))
		buildDisplay3d(4, 10, colorOptions, []string{string(wordDisplay)})
		buildDisplay3d(10, 4, colorPointingAt, []string{"Lettres déjà jouées : " + string(runesPlayed)})
		switch status {
		case ALREADYPLAYED:
			buildDisplay3d(11, 6, Orangered, []string{"Vous avez déjà joué cette lettre !"})
		case INCORRECTRUNE:
			buildDisplay3d(11, 6, Red, []string{"Lettre incorrecte : " + strconv.Itoa(len(hangman)-1-nbErrors) + " essais restants"})
		case CORRECTWORD:
			buildDisplay3d(11, 6, Limegreen, []string{"Félicitations, vous avez deviné le mot !"})
		case INCORRECTWORD:
			buildDisplay3d(11, 6, Red, []string{"Mot incorrect : " + strconv.Itoa(len(hangman)-1-nbErrors) + " essais restants"})
		default:
			break
		}
		if wordMode {
			buildDisplay3d(12, 4, colorOptions, []string{"Essayez de deviner le mot :"})
			buildDisplay3d(12, 32, colorPointingAt, []string{try})
		} else {
			buildDisplay3d(12, 4, colorOptions, []string{"Tapez une lettre pour deviner le mot"})
		}
		buildDisplay3d(13, 4, colorTitle, []string{"Utilisez les flèches (haut et bas) pour changer de mode"})
		showDisplay3d()
		if strings.Join(strings.Split(string(wordDisplay), " "), "") == strings.ToUpper(string(word)) {
			hasWon = true
			time.Sleep(time.Second * 2)
			break
		}
		if wordMode {
			var singleInput rune
			var enterPressed bool
			var charMode bool
			singleInput, enterPressed, charMode = wordInput()
			if !enterPressed && !charMode {
				if singleInput != 8 {
					try += string(singleInput)
				} else {
					if len(try) > 0 {
						try = string([]rune(try)[:len(try)-1])
					}
				}
			} else {
				if enterPressed {
					if checkWord(word, try) {
						score += nbRemainingLetters(wordDisplay) * 2
						status = CORRECTWORD
						revealWord(word, wordDisplay)
						try = ""
					} else {
						score -= nbRemainingLetters(wordDisplay) * 2
						nbErrors += 2
						status = INCORRECTWORD
						try = ""
					}
				}
				if charMode {
					wordMode = false
					try = ""
				}
			}
		} else {
			char, status, wordMode = input()
			displayWord(word, wordDisplay, char)
			if wordMode {
				status = CORRECTRUNE
			}
		}
		if nbErrors >= len(hangman)-1 {
			break
		}
	}
	endGame(hasWon)
}

func endGame(hasWon bool) {
	if hasWon {
		buildDisplay3d(4, 4, colorTitle, []string{"Félicitations,", "vous avez gagné !", "", "Le mot était : " + strings.ToUpper(string(word)), "", "Votre score est : " + strconv.Itoa(score), "", "", "Attendez quelques secondes pour revenir au menu principal"})
		showDisplay3d()
	} else {
		buildDisplay3d(6, 4, colorTitle, []string{"GAME OVER", "", "Le mot était : " + strings.ToUpper(string(word)), "", "", "Attendez quelques secondes pour revenir au menu principal"})
		buildDisplay3d(3, 52, colorTitle, strings.Split(hangman[len(hangman)-1], "\n"))
		showDisplay3d()
	}
	time.Sleep(time.Second * 10)
}

func clearGameData() {
	nbErrors = 0
	nbLettersFound = 0
	score = 0
	runesPlayed = append(runesPlayed[0:0])
}
