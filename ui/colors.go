package ui

/*
*
* @brief Return ANSI value of a selected color. Add bg prefix for background color.
black, red, green, yellow, blue, magneta, cyan, white,
brightBlack, brightRed, brightGreen, brightYellow, brightBlue, brightMagneta, brightCyan, brightWhite
*
* @param string
*
* @return
*/
func GetColor(colorName string) string {
	switch colorName {
	// Foreground colors
	case "black":
		return "\033[30m"
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "magenta":
		return "\033[35m"
	case "cyan":
		return "\033[36m"
	case "white":
		return "\033[37m"
	case "brightBlack":
		return "\033[90m"
	case "brightRed":
		return "\033[91m"
	case "brightGreen":
		return "\033[92m"
	case "brightYellow":
		return "\033[93m"
	case "brightBlue":
		return "\033[94m"
	case "brightMagenta":
		return "\033[95m"
	case "brightCyan":
		return "\033[96m"
	case "brightWhite":
		return "\033[97m"
	// Background colors
	case "bgBlack":
		return "\033[40m"
	case "bgRed":
		return "\033[41m"
	case "bgGreen":
		return "\033[42m"
	case "bgYellow":
		return "\033[43m"
	case "bgBlue":
		return "\033[44m"
	case "bgMagenta":
		return "\033[45m"
	case "bgCyan":
		return "\033[46m"
	case "bgWhite":
		return "\033[47m"
	case "bgBrightBlack":
		return "\033[100m"
	case "bgBrightRed":
		return "\033[101m"
	case "bgBrightGreen":
		return "\033[102m"
	case "bgBrightYellow":
		return "\033[103m"
	case "bgBrightBlue":
		return "\033[104m"
	case "bgBrightMagenta":
		return "\033[105m"
	case "bgBrightCyan":
		return "\033[106m"
	case "bgBrightWhite":
		return "\033[107m"
	default:
		return "\033[0m" // Default to reset if color not found
	}
}
