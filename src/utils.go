package ProjetHangman

import (
	"fmt"
	"github.com/mattn/go-tty"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var words []string
var word, runesPlayed []rune

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
		buildDisplay(0, colorTitle, []string{title})
		for i, option := range options {
			if cursorAt == i {
				buildDisplay(i+2, colorPointingAt, []string{"    " + cursor + "\t" + option})
			} else {
				buildDisplay(i+2, colorOptions, []string{"    " + "\t" + option})
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
				fmt.Println("Lettre correcte !")
				runesPlayed = append(runesPlayed, char)
			case INCORRECTRUNE:
				fmt.Println("Lettre incorrecte !")
				runesPlayed = append(runesPlayed, char)
			case ALREADYPLAYED:
				fmt.Println("Lettre déjà jouée !")
			}
		} else {
			fmt.Println("Il faut rentrer une lettre !")
		}
	}
}

func play() {
	clearTerminal()
	retreiveWords()
	word = []rune((words[rand.Intn(len(words)-1)]))
	var wordDisplay []rune
	for i := range word {
		wordDisplay = append(wordDisplay, '_')
		if i != len(word)-1 {
			wordDisplay = append(wordDisplay, ' ')
		}
	}
	fmt.Println(string(wordDisplay))
	input()
}
