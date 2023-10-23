package ProjetHangman

import (
	"regexp"
	"strings"
)

var display []string

var border = []string{
	"   _.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._   ",
	" ,'_.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._`. ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (                                                                     ) )",
	" ) )                                                                   ( ( ",
	"( (_.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._) )",
	" `._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._.-._,' ",
}

func buildDisplay(line, column int, color Color, content []string) {
	for i, str := range content {
		currentColumn := column
		if line+i > 13 {
			break
		}
		runesLineDisplay := []rune(display[line+2+i][4 : len(display[line+2+i])-4])
		colorRegExp := regexp.MustCompile("\\033\\[[0-9;]+m")
		colorCodes := colorRegExp.FindAllString(str, -1)
		colorIndexes := colorRegExp.FindAllStringIndex(str, -1)
		str = strings.Join(colorRegExp.Split(str, -1), "")
		prevColorCodes := colorRegExp.FindAllString(string(runesLineDisplay), -1)
		prevColorIndexes := colorRegExp.FindAllStringIndex(string(runesLineDisplay), -1)
		runesLineDisplay = []rune(strings.Join(colorRegExp.Split(string(runesLineDisplay), -1), ""))

		firstChar := true
		var firstColumn int
		for _, char := range str {
			if firstChar && char != ' ' {
				firstColumn = currentColumn
				firstChar = false
			}
			if char == '\t' {
				switch {
				case currentColumn < 14:
					currentColumn = 14
				case currentColumn < 22:
					currentColumn = 22
				case currentColumn < 30:
					currentColumn = 30
				case currentColumn < 40:
					currentColumn = 40
				case currentColumn < 50:
					currentColumn = 50
				case currentColumn < 60:
					currentColumn = 60
				case currentColumn > 60:
				}
			} else if char != ' ' {
				if len(runesLineDisplay) <= currentColumn {
					break
				} else {
					runesLineDisplay[currentColumn] = char
				}
			}
			currentColumn++
		}
		if colorCodes != nil && colorIndexes != nil {
			for i, _ := range colorCodes {
				runesLineDisplay = append(runesLineDisplay[:colorIndexes[i][0]+firstColumn], []rune(colorCodes[i]+string(runesLineDisplay[colorIndexes[i][0]+firstColumn:]))...)
			}
		}
		if prevColorCodes != nil && prevColorIndexes != nil {
			var indent int
			if firstColumn < prevColorIndexes[0][0] {
				indent = column - firstColumn
			}
			for i, _ := range prevColorCodes {
				runesLineDisplay = append(runesLineDisplay[:prevColorIndexes[i][0]+indent], []rune(prevColorCodes[i]+string(runesLineDisplay[prevColorIndexes[i][0]+indent:]))...)
			}
		}
		display[line+2+i] = display[line+2+i][:4] + colorCode(color) + string(runesLineDisplay) + CLEARCOLOR + display[line+2+i][len(display[line+2+i])-4:]
	}
}

func showDisplay() {
	clearTerminal()
	for _, str := range display {
		printColor(colorBorder, str[:len(str)-4], colorCode(colorBorder), str[len(str)-4:])
	}
	clearDisplay()
}

func clearDisplay() {
	display = append(display[0:0], border...)
}
