package main

import (
	"flag"
	"fmt"
	//	"github.com/aeden/traceroute"
	"github.com/nsf/termbox-go"
)

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for n := 0; n < len(runes); n += 1 {
		termbox.SetCell(x+n, y, runes[n], color, backgroundColor)
	}
}

func drawLineFull(x, y int, str string, lineAttr, backAttr termbox.Attribute) {
	color := lineAttr
	backgroundColor := backAttr
	runes := []rune(str)

	for n := 0; n < len(runes); n += 1 {
		termbox.SetCell(x+n, y, runes[n], color, backgroundColor)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	defer termbox.Close()

	text := make([]string, 0, 30)
	maxX, maxY := termbox.Size()

	first := flag.Arg(0)
	huff := first

	cursX := 1
	cursY := maxY - 2
	termbox.SetCursor(cursX+1, 2)

	var cmdLine string
	for i := 0; i < maxX; i++ {
		cmdLine = cmdLine + " "
	}
	drawLineFull(0, cursY, cmdLine, termbox.ColorDefault, termbox.ColorRed)
	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyBackspace:
				if cursX > 1 {
					cursX--
					termbox.SetCursor(cursX+1, cursY)
					drawLine(cursX+1, cursY, " ")

					text = text[:len(text)-1]
					termbox.Flush()
				}
			case termbox.KeyEnter:
				x := 1
				huff = ""
				fill(x, cursY, maxX-x, 2, termbox.Cell{Ch: ' '})
				drawLineFull(0, cursY, cmdLine, termbox.ColorDefault, termbox.ColorRed)
				for _, s := range text {
					drawLineFull(x+1, cursY+1, s, termbox.ColorRed, termbox.ColorDefault)
					huff = huff + s
					x++
				}
				text = make([]string, 0, 30)
				cursX = 1
				termbox.SetCursor(cursX+1, cursY)
				termbox.Flush()
			default:
				if cursX < maxX-1 {
					cursX++
					termbox.SetCursor(cursX+1, cursY)
					drawLineFull(cursX, cursY, fmt.Sprintf("%c", ev.Ch), termbox.ColorDefault, termbox.ColorRed)
					termbox.Flush()
					text = append(text, fmt.Sprintf("%c", ev.Ch))
				}
			}
		}
	}
}
