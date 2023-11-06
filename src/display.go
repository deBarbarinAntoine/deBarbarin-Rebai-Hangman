package ProjetHangman

import (
	"fmt"
	"regexp"
	"strings"
)

var display3d [18][75]string

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

var scoreDisplay = []string{
	"    │       NOM       │ SCORE │   DIFFICULTÉ   │   DICTIONNAIRE    ",
	"────┼─────────────────┼───────┼────────────────┼───────────────────",
	"  1 │                 │       │                │                   ",
	"  2 │                 │       │                │                   ",
	"  3 │                 │       │                │                   ",
	"  4 │                 │       │                │                   ",
	"  5 │                 │       │                │                   ",
	"  6 │                 │       │                │                   ",
	"  7 │                 │       │                │                   ",
	"  8 │                 │       │                │                   ",
	"  9 │                 │       │                │                   ",
	" 10 │                 │       │                │                   ",
	"              [Tapez sur Entrée pour revenir au menu]              ",
}

// buildDisplay3d adds the content to display3d with color codes in the line and column specified in the parameters (for each character that appears on screen, there is a string that contains the color code, the character, and the CLEARCOLOR to reset the color code. If the content contains a color code, it will change the color put in the parameters from the corresponding column onwards.
func buildDisplay3d(line, column int, color Color, content []string) {
	currentColorCode := colorCode(color)
	line += 2
	column += 4
	for i, str := range content {
		currentColumn := column
		if line+i > 15 {
			break
		}
		colorRegExp := regexp.MustCompile("\\033\\[[0-9;]+m")
		colorCodes := colorRegExp.FindAllString(str, -1)
		var colorIndexes []int
		for i := 0; i < len(colorCodes); i++ {
			colorIndexes = append(colorIndexes, colorRegExp.FindStringIndex(str)[0])
			str = strings.Join(colorRegExp.Split(str, 2), "")
		}

		for k, char := range str {
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
				if currentColumn >= 72 {
					break
				} else {
					if colorCodes != nil && colorIndexes != nil {
						for j, colorColumn := range colorIndexes {
							if k >= colorColumn {
								currentColorCode = colorCodes[j]
							}
						}
					}
					display3d[line+i][currentColumn] = currentColorCode + string(char) + CLEARCOLOR
				}
			}
			currentColumn++
		}
	}
}

// showDisplay3d prints display3d in the terminal and resets it afterwards.
func showDisplay3d() {
	clearTerminal()
	for _, line := range display3d {
		for _, char := range line {
			fmt.Print(char)
		}
		fmt.Println()
	}
	clearDisplay3d()
}

// clearDisplay3d clears the variable display3d and adds the border.
func clearDisplay3d() {
	display3d = [18][75]string{}
	for i, str := range border {
		for j, char := range str {
			display3d[i][j] = colorCode(colorBorder) + string(char) + CLEARCOLOR
		}
	}
}
