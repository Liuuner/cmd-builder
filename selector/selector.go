package selector

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"errors"
	"fmt"
	"github.com/liuuner/selto/colors"
	"io"
	"os"
)

var col = colors.CreateColors(true)

type Item struct {
	Id      string
	Display string
	Color   colors.Formatter
}

type Selector struct {
	writer          io.Writer
	cursorString    string
	promptString    string
	completedString string
	failedString    string
	hasPrompt       bool
	hasSummary      bool
	title           string
	items           []Item
	cursorPos       int
}

type Option func(*Selector)

func New(items []Item, title string, opts ...Option) *Selector {
	s := &Selector{
		writer:          os.Stdout,
		title:           title,
		items:           items,
		cursorString:    "❯",
		promptString:    "?",
		completedString: "✔",
		failedString:    "✖",
		hasPrompt:       true,
		hasSummary:      true,
		cursorPos:       0,
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Selector) Open() (item Item, err error) {
	cursorHide()

	if s.hasPrompt {
		fmt.Printf("%s %s: %s\n",
			col.Cyan(s.promptString),
			s.title,
			col.Gray("› - Use arrow-keys. Return to submit."),
		)
	}

	s.render(false)

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape, keys.Backspace:
			return true, errors.New("canceled") // Return true to stop listener
		case keys.Up:
			s.cursorPos--
			s.keepPosInBoundaries()
			s.render(true)
		case keys.Down:
			s.cursorPos++
			s.keepPosInBoundaries()
			s.render(true)
		case keys.Enter:
			return true, nil
		}

		return false, nil
	})

	s.clear()

	if s.hasSummary && err != nil {
		fmt.Printf("%s %s: %s\n", col.Red(s.failedString), s.title, col.Gray("› ", err))
		return Item{}, err
	}
	selectedItem := s.items[s.cursorPos]

	if s.hasSummary {
		if selectedItem.Color != nil {
			fmt.Printf("%s %s: %s %s\n", col.Green(s.completedString), s.title, col.Gray("›"), selectedItem.Color(selectedItem.Display))
		} else {
			fmt.Printf("%s %s: %s %s\n", col.Green(s.completedString), s.title, col.Gray("›"), selectedItem.Display)
		}
	}

	cursorShow()

	return selectedItem, nil
}

func (s *Selector) keepPosInBoundaries() {
	s.cursorPos = (s.cursorPos + len(s.items)) % len(s.items)
}

func (s *Selector) render(rerender bool) {
	if rerender {
		// Move cursor to top
		s.writer.Write([]byte(fmt.Sprintf("\033[%dA", len(s.items)-1)))
	}

	for index, item := range s.items {
		var newline = "\n"
		if index == len(s.items)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := item.Display
		if item.Color != nil {
			menuItemText = item.Color(item.Display)
		}
		cursor := "   "
		if index == s.cursorPos { // for color or other effects
			cursor = col.Cyan(s.cursorString, "  ")
			menuItemText = col.Underline(menuItemText)
		}

		s.writer.Write([]byte(fmt.Sprintf("\r%s %s%s", cursor, menuItemText, newline)))
	}
}

func (s *Selector) clear() {
	clearLine()
	for i := 0; i < len(s.items)-1; i++ {
		cursorUp()
		clearLine()
	}

	if s.hasPrompt {
		cursorUp()
		clearLine()
	}
}

func cursorUp() {
	fmt.Print("\033[F") // ANSI escape code to move cursor up
}

func clearLine() {
	fmt.Print("\u001b[2K") // ANSI escape code to clear the line
}

func cursorShow() {
	fmt.Printf("\033[?25h") // Show Cursor
}

func cursorHide() {
	fmt.Printf("\033[?25l") // Hide Cursor
}

// Options
func WithWriter(writer io.Writer) Option {
	return func(s *Selector) {
		s.writer = writer
	}
}

func WithCursorString(cursorString string) Option {
	return func(s *Selector) {
		s.cursorString = cursorString
	}
}

func WithPrompt(b bool) Option {
	return func(s *Selector) {
		s.hasPrompt = b
	}
}

func WithSummary(b bool) Option {
	return func(s *Selector) {
		s.hasSummary = b
	}
}

func WithInitialCursorPos(i int) Option {
	return func(s *Selector) {
		s.cursorPos = i
	}
}
