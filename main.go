package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v3"
)

func main() {
	filename := flag.String("f", "", "file to open")
	showHelp := flag.Bool("help", false, "show help")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("initializing screen: %v", err)
	}
	defer screen.Fini()

	width, height := screen.Size()
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	statusStyle := style.Background(tcell.ColorDarkCyan)

	draw := func() {
		screen.Clear()
		width, height = screen.Size()
		for x := 0; x < width; x++ {
			screen.SetContent(x, height-1, ' ', nil, statusStyle)
		}
		statusMsg := fmt.Sprintf(" et — %s | Ctrl+Q quit", *filename)
		if *filename == "" {
			statusMsg = " et — <new file> | Ctrl+Q quit"
		}
		for i, ch := range statusMsg {
			if i >= width {
				break
			}
			screen.SetContent(i, height-1, ch, nil, statusStyle)
		}
		screen.Show()
	}

	draw()

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			draw()
		case *tcell.EventKey:
			if e.Key() == tcell.KeyCtrlQ || e.Key() == tcell.KeyEscape {
				return
			}
		}
	}
}
