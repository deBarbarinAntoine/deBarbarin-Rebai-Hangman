package main

import (
	"fmt"
	"github.com/mattn/go-tty"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var words []string
var word, runesPlayed []rune

const (
	ALREADYPLAYED = 0
	CORRECTRUNE   = 1
	INCORRECTRUNE = 2
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
		clearTerminal()
		fmt.Println("\033[38;2;30;144;255m", title, "\033[0m")
		for i, option := range options {
			if cursorAt == i {
				fmt.Println("\033[38;2;0;128;128m", cursor, "\t", option, "\033[0m")
			} else {
				fmt.Println("\033[38;2;30;144;255m", "\t", option, "\033[0m")
			}
		}
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

func principalMenu() {
	for {
		switch createVerticalMenu(0, "-->", "------- MENU PRINCIPAL -------", "Nouvelle partie", "Paramètres", "Quitter") {
		case "Nouvelle partie":
			play()
		case "Paramètres":
			clearTerminal()
			fmt.Println("Paramètres...")
			time.Sleep(time.Second * 2)
		case "Quitter":
			clearTerminal()
			os.Exit(1)
		}
	}
}

func retreiveWords() {
	content, err := os.ReadFile("Files/Dictionaries/ods5.txt")
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
	var display []rune
	for i := range word {
		display = append(display, '_')
		if i != len(word)-1 {
			display = append(display, ' ')
		}
	}
	fmt.Println(string(display))
	input()
}
