package colors

import "strings"

// todo is color supported (--no-color usw)
// https://github.com/alexeyraspopov/picocolors/blob/main/picocolors.js

type Formatter func(input string) string

func formatter(open, close, replace string) Formatter {
	if replace == "" {
		replace = open
	}
	return func(input string) string {
		index := strings.Index(input, close)
		if index != -1 {
			return open + replaceClose(input, close, replace, index) + close
		}
		return open + input + close
	}
}

func replaceClose(str, close, replace string, index int) string {
	result := ""
	cursor := 0
	for index != -1 {
		result += str[cursor:index] + replace
		cursor = index + len(close)
		index = strings.Index(str[cursor:], close)
		if index != -1 {
			index += cursor
		}
	}
	return result + str[cursor:]
}

type Colors struct {
	IsColorSupported bool
	Reset            Formatter
	Bold             Formatter
	Dim              Formatter
	Italic           Formatter
	Underline        Formatter
	Inverse          Formatter
	Hidden           Formatter
	Strikethrough    Formatter
	Black            Formatter
	Red              Formatter
	Green            Formatter
	Yellow           Formatter
	Blue             Formatter
	Magenta          Formatter
	Cyan             Formatter
	White            Formatter
	Gray             Formatter
	BgBlack          Formatter
	BgRed            Formatter
	BgGreen          Formatter
	BgYellow         Formatter
	BgBlue           Formatter
	BgMagenta        Formatter
	BgCyan           Formatter
	BgWhite          Formatter
	BlackBright      Formatter
	RedBright        Formatter
	GreenBright      Formatter
	YellowBright     Formatter
	BlueBright       Formatter
	MagentaBright    Formatter
	CyanBright       Formatter
	WhiteBright      Formatter
	BgBlackBright    Formatter
	BgRedBright      Formatter
	BgGreenBright    Formatter
	BgYellowBright   Formatter
	BgBlueBright     Formatter
	BgMagentaBright  Formatter
	BgCyanBright     Formatter
	BgWhiteBright    Formatter
}

func CreateColors(enabled bool) Colors {
	init := func(open, close string, replace ...string) Formatter {
		r := open
		if len(replace) > 0 {
			r = replace[0]
		}
		if enabled {
			return formatter(open, close, r)
		}
		return func(input string) string {
			return input
		}
	}

	return Colors{
		IsColorSupported: enabled,
		Reset:            init("\x1b[0m", "\x1b[0m"),
		Bold:             init("\x1b[1m", "\x1b[22m", "\x1b[22m\x1b[1m"),
		Dim:              init("\x1b[2m", "\x1b[22m", "\x1b[22m\x1b[2m"),
		Italic:           init("\x1b[3m", "\x1b[23m"),
		Underline:        init("\x1b[4m", "\x1b[24m"),
		Inverse:          init("\x1b[7m", "\x1b[27m"),
		Hidden:           init("\x1b[8m", "\x1b[28m"),
		Strikethrough:    init("\x1b[9m", "\x1b[29m"),

		Black:   init("\x1b[30m", "\x1b[39m"),
		Red:     init("\x1b[31m", "\x1b[39m"),
		Green:   init("\x1b[32m", "\x1b[39m"),
		Yellow:  init("\x1b[33m", "\x1b[39m"),
		Blue:    init("\x1b[34m", "\x1b[39m"),
		Magenta: init("\x1b[35m", "\x1b[39m"),
		Cyan:    init("\x1b[36m", "\x1b[39m"),
		White:   init("\x1b[37m", "\x1b[39m"),
		Gray:    init("\x1b[90m", "\x1b[39m"),

		BgBlack:         init("\x1b[40m", "\x1b[49m"),
		BgRed:           init("\x1b[41m", "\x1b[49m"),
		BgGreen:         init("\x1b[42m", "\x1b[49m"),
		BgYellow:        init("\x1b[43m", "\x1b[49m"),
		BgBlue:          init("\x1b[44m", "\x1b[49m"),
		BgMagenta:       init("\x1b[45m", "\x1b[49m"),
		BgCyan:          init("\x1b[46m", "\x1b[49m"),
		BgWhite:         init("\x1b[47m", "\x1b[49m"),
		BlackBright:     init("\x1b[90m", "\x1b[39m"),
		RedBright:       init("\x1b[91m", "\x1b[39m"),
		GreenBright:     init("\x1b[92m", "\x1b[39m"),
		YellowBright:    init("\x1b[93m", "\x1b[39m"),
		BlueBright:      init("\x1b[94m", "\x1b[39m"),
		MagentaBright:   init("\x1b[95m", "\x1b[39m"),
		CyanBright:      init("\x1b[96m", "\x1b[39m"),
		WhiteBright:     init("\x1b[97m", "\x1b[39m"),
		BgBlackBright:   init("\x1b[100m", "\x1b[49m"),
		BgRedBright:     init("\x1b[101m", "\x1b[49m"),
		BgGreenBright:   init("\x1b[102m", "\x1b[49m"),
		BgYellowBright:  init("\x1b[103m", "\x1b[49m"),
		BgBlueBright:    init("\x1b[104m", "\x1b[49m"),
		BgMagentaBright: init("\x1b[105m", "\x1b[49m"),
		BgCyanBright:    init("\x1b[106m", "\x1b[49m"),
		BgWhiteBright:   init("\x1b[107m", "\x1b[49m"),
	}
}
