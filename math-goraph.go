package main

import (
	"flag"
	"fmt"
	//	"github.com/aeden/traceroute"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	INIT_CURS_X = 5
	YEQ_LINE    = "y ="
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

func initFill(maxX, cursY int) {
	maxY := cursY - 2
	midX := maxX / 2
	midY := maxY / 2
	Xpunc := midX / 5
	Ypunc := midY / 5

	fill(0, 0, maxX, maxY, termbox.Cell{Ch: ' '})
	fill(0, midY, maxX, 1, termbox.Cell{Ch: '-'})
	fill(midX, 0, 2, maxY, termbox.Cell{Ch: '|'})
	drawLine(midX+2, midY+1, "O")

	var YpunCounter int
	for i := 1; i <= Ypunc; i++ {
		YpunCounter += 5
		drawLine(midX+2, midY-YpunCounter, fmt.Sprint(YpunCounter))
		drawLine(midX+2, midY+YpunCounter, fmt.Sprintf("-%v", YpunCounter))
	}
	var XpunCounter int
	for i := 1; i <= Xpunc; i++ {
		XpunCounter += 5
		drawLine(midX+(XpunCounter*2), midY+1, fmt.Sprint(XpunCounter))
		drawLine(midX-(XpunCounter*2), midY+1, fmt.Sprintf("-%v", XpunCounter))
	}
}

func selectMode(mode string, maxX, cursY int, kill chan struct{}) {
	switch mode {
	case "c":
		drawCubicLoop(maxX, cursY, kill)
	case "e":
		//drawExponentialLoop(maxX, cursY, kill)
	case "l":
		//drawLogarithmLoop(maxX, cursY, kill)
	case "t":
		//drawTrigonometricLoop(maxX, cursY, kill)
	}
}

func drawCubicLoop(maxX, cursY int, kill chan struct{}) {
	maxY := cursY - 1
	midX := maxX / 2
	midY := maxY / 2
	var firstX int
	if midX < midY {
		firstX = maxY - (maxX * 2)
	} else {
		firstX = (maxX / 2) - maxY
	}
	if (maxX % 2) == 0 {
		firstX++
	}

	/*if y = x*/
	var i int
	for {
		select {
		case <-kill:
			return
		default:

			if maxY-i < 0 {
				i = 0
			}
			drawLineFull(firstX+i*2, maxY-i, "||", termbox.ColorWhite, termbox.Attribute(18+1))
			termbox.Flush()
			time.Sleep(100 * time.Millisecond)
			drawLineFull(firstX+i*2, maxY-i, "  ", termbox.ColorGreen, termbox.Attribute(100+1))
			termbox.Flush()
			i++
		}
	}
}

func drawCoodinate(maxX, maxY, cursX int) {
	cursY := maxY - 2
	var cmdLine string
	coordinate := fmt.Sprintf("maxX: %v maxY: %v ", (maxX / 2), maxY)
	for i := 0; i < maxX-len(coordinate); i++ {
		cmdLine = cmdLine + " "
	}
	cmdLine = cmdLine + coordinate
	fill(maxX-len(coordinate), cursY, len(coordinate), 1, termbox.Cell{Ch: ' '})
	drawLineFull(0, cursY, cmdLine, termbox.ColorDefault, termbox.ColorRed)
	drawLineFull(1, cursY, YEQ_LINE, termbox.ColorYellow, termbox.ColorRed)
	drawLineFull(1, cursY, YEQ_LINE, termbox.ColorYellow, termbox.ColorRed)
	termbox.SetCursor(cursX+1, cursY)
	termbox.Flush()
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	defer termbox.Close()

	modeMap := map[string]string{
		"MODE_CUBIC":         "c", // 1~3次間数
		"MODE_EXPONENTIAL":   "e", // 指数関数
		"MODE_LOGARITHM":     "l", // 対数関数
		"MODE_TRIGONOMETRIC": "t", // 三角関数
	}

	modeView := map[string]string{
		"c": "[mode: cubic]",
		"e": "[mode: exponential]",
		"l": "[mode: logarithm]",
		"t": "[mode: trigonometric]",
	}

	opeMap := map[string]string{
		"c": "x",
		"e": "th power of x",
		"l": " base's logX",
		"t": "of sin?",
	}

	text := make([]string, 0, 30)
	maxX, maxY := termbox.Size()

	first := flag.Arg(0)
	huff := first
	kill := make(chan struct{}, 0)

	cursX := INIT_CURS_X
	mode := modeMap["MODE_CUBIC"]
	ope := opeMap[mode]
	cursY := maxY - 2
	termbox.SetCursor(cursX+1, cursY)

	drawCoodinate(maxX, maxY, cursX)
	initFill(maxX, maxY)
	drawLineFull(cursX+2, cursY, ope, termbox.ColorDefault, termbox.ColorRed)
	drawLineFull(1, cursY-1, "[mode: cubic]", termbox.Attribute(21+1), termbox.ColorDefault)
	termbox.Flush()

	go selectMode(mode, maxX, cursY, kill)

	var tabCount int
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			maxX, maxY = termbox.Size()
			initFill(maxX, maxY)
			cursX = INIT_CURS_X
			drawCoodinate(maxX, maxY, cursX)
			kill <- struct{}{}
			go selectMode(mode, maxX, cursY, kill)
			termbox.Flush()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeyEnter:
				kill <- struct{}{}
				x := INIT_CURS_X
				huff = ""
				fill(x, cursY, maxX-x, 2, termbox.Cell{Ch: ' '})
				for _, s := range text {
					drawLineFull(x+1, cursY+1, s, termbox.ColorRed, termbox.ColorDefault)
					huff = huff + s
					x++
				}
				text = make([]string, 0, 30)
				cursX = INIT_CURS_X
				drawCoodinate(maxX, maxY, cursX)
				initFill(maxX, maxY)
				drawLineFull(1, cursY+1, YEQ_LINE, termbox.ColorYellow, termbox.ColorDefault)
				drawLineFull(cursX+2, cursY, ope, termbox.ColorDefault, termbox.ColorRed)
				termbox.Flush()
				go selectMode(mode, maxX, cursY, kill)
			case termbox.KeyTab:
				tabCount++
				fill(1, cursY-1, len(modeView[mode]), 1, termbox.Cell{Ch: ' '})
				t := tabCount % 4
				switch t {
				case 0:
					mode = modeMap["MODE_CUBIC"]
				case 1:
					mode = modeMap["MODE_EXPONENTIAL"]
				case 2:
					mode = modeMap["MODE_LOGARITHM"]
				case 3:
					mode = modeMap["MODE_TRIGONOMETRIC"]
				}
				drawLineFull(1, cursY-1, modeView[mode], termbox.Attribute(21+1), termbox.ColorDefault)
				termbox.Flush()
			case termbox.KeyBackspace:
				if cursX > INIT_CURS_X {
					drawLineFull(cursX+2, cursY, " ", termbox.ColorDefault, termbox.ColorRed)
					cursX--
					termbox.SetCursor(cursX+1, cursY)
					drawLineFull(cursX+1, cursY, " ", termbox.ColorDefault, termbox.ColorRed)
					drawLineFull(cursX+2, cursY, ope, termbox.ColorDefault, termbox.ColorRed)
					text = text[:len(text)-1]
					termbox.Flush()
				}
			default:
				if cursX < maxX-1 {
					drawLineFull(cursX+2, cursY, " ", termbox.ColorDefault, termbox.ColorRed)
					cursX++
					termbox.SetCursor(cursX+1, cursY)
					drawLineFull(cursX, cursY, fmt.Sprintf("%c", ev.Ch), termbox.ColorDefault, termbox.ColorRed)
					drawLineFull(cursX+2, cursY, ope, termbox.ColorDefault, termbox.ColorRed)
					termbox.Flush()
					text = append(text, fmt.Sprintf("%c", ev.Ch))
				}
			}
		}
	}
}
