package selector

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"context"
	"fmt"
	"github.com/liuuner/temp-container/colors"
	"io"
	"os"
	"sync"
)

var col = colors.CreateColors(true)

type Item struct {
	Name    string
	Display string
	Color   colors.Formatter
}

type Selector struct {
	writer     io.Writer
	cursorPos  int
	cursorChar rune
	items      []Item
	cancelFunc context.CancelFunc
	doneCh     chan struct{}
	lock       sync.RWMutex
	hasTitle   bool
}

type Config struct {
	Writer     io.Writer
	CursorChar rune
}

func New(items []Item, cfg Config) *Selector {
	s := &Selector{
		writer:     os.Stderr,
		cursorPos:  0,
		cursorChar: '>',
		items:      items,
		hasTitle:   true,
	}

	if cfg.Writer != nil {
		s.writer = cfg.Writer
	}

	if cfg.CursorChar != 0 {
		s.cursorChar = cfg.CursorChar
	}

	return s
}

func (s *Selector) Open() Item {
	/*defer func() {
		// show the cursor
		fmt.Printf("\033[?25h")
	}()*/

	// remove the cursor
	fmt.Printf("\033[?25l")
	s.lock.Lock()

	if s.hasTitle {
		fmt.Printf("%s Select a framework: %s\n",
			col.Cyan("?"),
			col.Gray("› - Use arrow-keys. Return to submit."),
		)
	}

	s.render(false)

	//ctx, cancel := context.WithCancel(context.Background())

	// Stop keyboard listener on Escape key press or CTRL+C.
	// Exit application on "q" key press.
	// Print every rune key press.
	// Print every other key press.
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC, keys.Escape:
			return true, nil // Return true to stop listener
		case keys.Up:
			s.cursorPos--
			s.keepPosInBoudaries()
			s.render(true)
		case keys.Down:
			s.cursorPos++
			s.keepPosInBoudaries()
			s.render(true)

		case keys.Enter:
			return true, nil
		}

		return false, nil

		/*case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
			if key.String() == "q" { // Check if key is "q"
				fmt.Println("\rQuitting application")
				os.Exit(0) // Exit application
			}
			fmt.Printf("\rYou pressed the rune key: %s\n", key)
		default:
			fmt.Printf("\rYou pressed: %s\n", key)
		}*/

		// Return false to continue listening
	})
	fmt.Printf("\033[?25h")
	s.clear()
	//fmt.Print("#")
	selectedItem := s.items[s.cursorPos]

	if s.hasTitle {
		fmt.Printf("%s Select a framework: %s", col.Green("✔"), col.Gray("› "+selectedItem.Display))
	}

	return selectedItem
}

func (s *Selector) keepPosInBoudaries() {
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
		cursor := "  "
		if index == s.cursorPos { // for color or other effects
			cursor = col.Cyan("> ")
			menuItemText = col.Bold(menuItemText)
		}

		s.writer.Write([]byte(fmt.Sprintf("\r%s %s%s", cursor, menuItemText, newline)))
	}
}

func (s *Selector) clear() {
	fmt.Print("\u001b[2K") // ANSI escape code to clear the line
	for i := 0; i < len(s.items)-1; i++ {
		fmt.Print("\033[F")    // ANSI escape code to move cursor up
		fmt.Print("\u001b[2K") // ANSI escape code to clear the line
	}

	if s.hasTitle {
		fmt.Print("\033[F")    // ANSI escape code to move cursor up
		fmt.Print("\u001b[2K") // ANSI escape code to clear the line
	}
}

func (s *Selector) Close() {
	if !s.isOpen() {
		return
	}
	s.cancelFunc()

	s.lock.Lock()
	defer s.lock.Unlock()
}

func (s *Selector) isOpen() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.doneCh != nil
}
