package ProjetHangman

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

func buildDisplay(line int, color Color, content []string) {
	for i, str := range content {
		var column int
		if line+i > 13 {
			break
		}
		runesLineDisplay := []rune(display[line+2+i][4 : len(display[line+2+i])-4])
		for _, char := range str {
			if char == '\t' {
				switch {
				case column < 14:
					column = 14
				case column < 22:
					column = 22
				case column < 30:
					column = 30
				case column < 40:
					column = 40
				case column < 50:
					column = 50
				case column < 60:
					column = 60
				case column > 60:
				}
			} else if char != ' ' {
				if len(runesLineDisplay) <= column {
					break
				} else {
					runesLineDisplay[column] = char
				}
			}
			column++
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
	display = make([]string, 0)
	for _, str := range border {
		display = append(display, str)
	}
}
