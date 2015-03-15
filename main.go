package main

import (
	"github.com/nsf/termbox-go"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

var text *TextBuffer
var stdErr *log.Logger

type TextBuffer struct {
	cursor *Cursor
	runes  [][]rune
	width  int
	height int
}

func NewTextBuffer() *TextBuffer {
	w, h, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	stdErr.Printf("Dimensions: (%d, %d)", w, h)
	return &TextBuffer{
		cursor: &Cursor{},
		width:  w,
		height: h,
	}
}

func (text *TextBuffer) Draw() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
	for i := 0; i < len(text.runes) && i < text.height; i++ {
		row := text.runes[i]
		for j := 0; j < len(row) && j < text.width; j++ {
			termbox.SetCell(j, i, row[j], termbox.ColorWhite, termbox.ColorDefault)
		}
	}
	text.cursor.Draw()
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

func (text *TextBuffer) InsertRune(ch rune) {
	for len(text.runes)-1 < text.cursor.y {
		text.runes = append(text.runes, []rune{})
	}
	row := text.runes[text.cursor.y]
	row = append(row, ch)
	text.runes[text.cursor.y] = row
	text.cursor.x++
}

func (text *TextBuffer) RemoveRune() {
	if len(text.runes) == 0 {
		return
	}
	row := text.runes[text.cursor.y]
	if len(row) == 0 {
		if text.cursor.y == 0 {
			return
		}
		rowAbove := text.runes[text.cursor.y-1]
		text.cursor.x = len(rowAbove)
		text.cursor.y--
		text.runes = text.runes[:text.cursor.y+1]
		return
	}
	row = row[:len(row)-1]
	text.runes[text.cursor.y] = row
	text.cursor.x--
}

func (text *TextBuffer) InsertNewLine() {
	if len(text.runes)-1 < text.cursor.y {
		text.runes = append(text.runes, []rune{})
	} else {

	}
	text.cursor.x = 0
	text.cursor.y++
}

type Cursor struct {
	x int
	y int
}

func (c *Cursor) Draw() {
	termbox.SetCursor(c.x, c.y)
}

func init() {
	stdErr = log.New(os.Stderr, "", 0)
	stdErr.Println("\nInitializing...")
	text = NewTextBuffer()
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	text.Draw()

	for {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			stdErr.Println("Key event")
			switch event.Key {
			case termbox.KeyCtrlC:
				return
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				text.RemoveRune()
				text.Draw()
			case termbox.KeyEnter:
				text.InsertNewLine()
				text.Draw()
			default:
				text.InsertRune(event.Ch)
				text.Draw()
			}

		case termbox.EventResize:
			text.width = event.Width
			text.height = event.Height

		case termbox.EventError:
			panic(event.Err)

		case termbox.EventInterrupt:
			return

		}
	}
}
