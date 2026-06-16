package main

import (
	"et/config"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v3"
)

func draw(s tcell.Screen, cfg *config.Config, filename string) {
	style := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	statusStyle := style.Background(cfg.Colors.StatusBar.Color)
	s.Clear()
	w, h := s.Size()
	for x := range w {
		s.SetContent(x, h-1, ' ', nil, statusStyle)
	}
	statusMsg := fmt.Sprintf(" et — %s | Ctrl+Q quit", filename)
	if filename == "" {
		statusMsg = " et — <new file> | Ctrl+Q quit"
	}
	for i, ch := range statusMsg {
		if i >= w {
			break
		}
		s.SetContent(i, h-1, ch, nil, statusStyle)
	}
	s.Show()
}

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

	cfg := config.NewDefault()
	draw(screen, cfg, *filename)

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			draw(screen, cfg, *filename)
		case *tcell.EventKey:
			if e.Key() == tcell.KeyCtrlQ || e.Key() == tcell.KeyEscape {
				return
			}
		}
	}
}
