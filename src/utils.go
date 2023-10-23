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
var nbLettersFound int
var nbErrors int
var hangman []string

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

	CHANGECOLORBORDER           = 3
	CHANGECOLORTITLE            = 4
	CHANGECOLOROPTIONS          = 5
	CHANGECOLOROPTIONPOINTINGAT = 6
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
		clearDisplay()
		buildDisplay(0, 0, colorTitle, []string{title})
		for i, option := range options {
			if cursorAt == i {
				buildDisplay(i+2, 4, colorPointingAt, []string{cursor + "\t" + option})
			} else {
				buildDisplay(i+2, 4, colorOptions, []string{"\t" + option})
			}
		}
		showDisplay()
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
			clearDisplay()
			return options[cursorAt]
		}
	}
}

func PrincipalMenu() {
	clearTerminal()
	for {
		switch createVerticalMenu(0, "-->", "------------------------- MENU PRINCIPAL -------------------------", "Nouvelle partie", "Paramètres", "Quitter") {
		case "Nouvelle partie":
			play()
		case "Paramètres":
			parameters()
		case "Quitter":
			clearDisplay()
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
	content, err := os.ReadFile("../Files/hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	var line int
	var str string
	for _, char := range content {
		if line == 8 {
			hangman = append(hangman, str)
			str = ""
			line = 0
		}
		str += string(char)
		if char == '\n' {
			line++
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

func input() rune {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		char, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			switch checkRune(char) {
			case CORRECTRUNE:
				//fmt.Println("Lettre correcte !")
				runesPlayed = append(runesPlayed, char)
				nbLettersFound++
				return char
			case INCORRECTRUNE:
				//fmt.Println("Lettre incorrecte !")
				runesPlayed = append(runesPlayed, char)
				nbErrors++
				nbLettersFound--
				return '\n'
			case ALREADYPLAYED:
				//fmt.Println("Lettre déjà jouée !")
			}
		} else {
			//fmt.Println("Il faut rentrer une lettre !")
		}
	}
}

func play() {
	var hasWon bool
	clearTerminal()
	retreiveWords()
	retreiveHangman()
	clearGameData()
	word = []rune((words[rand.Intn(len(words)-1)]))
	wordDisplay := []rune(strings.Repeat("_ ", len(word)))
	for {
		buildDisplay(2, 4, colorTitle, []string{"Score : " + colorCode(colorPointingAt) + strconv.Itoa(nbLettersFound)})
		buildDisplay(3, 50, colorTitle, strings.Split(hangman[nbErrors], "\n"))
		buildDisplay(4, 10, colorOptions, []string{string(wordDisplay)})
		buildDisplay(10, 4, colorPointingAt, []string{"Lettres déjà jouées : " + string(runesPlayed)})
		buildDisplay(12, 4, colorOptions, []string{"Tapez une lettre pour deviner le mot"})
		//buildDisplay(12, 4, colorOptions, []string{"Essayez de deviner le mot"})
		buildDisplay(13, 4, colorTitle, []string{"Utilisez les flèches (haut et bas) pour changer de mode"})
		showDisplay()
		if strings.Join(strings.Split(string(wordDisplay), " "), "") == strings.ToUpper(string(word)) {
			hasWon = true
			time.Sleep(time.Second * 2)
			break
		}
		if nbErrors >= len(hangman) {
			break
		}
		displayWord(word, wordDisplay, input())
	}
	endGame(hasWon)
}

func endGame(hasWon bool) {
	if hasWon {
		buildDisplay(4, 4, colorTitle, []string{"Félicitations,", "vous avez gagné !", "", "Le mot était : " + strings.ToUpper(string(word)), "", "Votre score est : " + strconv.Itoa(nbLettersFound), "", "", "Attendez quelques secondes pour revenir au menu principal"})
		showDisplay()
	} else {
		buildDisplay(6, 4, colorTitle, []string{"GAME OVER", "", "Le mot était : " + strings.ToUpper(string(word)), "", "", "Attendez quelques secondes pour revenir au menu principal"})
		showDisplay()
	}
	time.Sleep(time.Second * 10)
}

func clearGameData() {
	nbErrors = 0
	nbLettersFound = 0
	runesPlayed = append(runesPlayed[0:0])
}
