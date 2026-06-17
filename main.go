package main

import (
	"et/config"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v3"
)

const logFileName = "et.log"
const version = "0.0.1"

func draw(s tcell.Screen, cfg *config.Config, filename string) {
	style := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	statusStyle := style.Background(cfg.Colors.StatusBar.Color)
	s.Clear()
	w, h := s.Size()
	for x := range w {
		s.SetContent(x, h-1, ' ', nil, statusStyle)
	}
	quitKeyBindsString := cfg.GetQuitKeyBindingsAsStr()
	fnameStr := filename
	if filename == "" {
		fnameStr = "<new file>"
	}
	statusMsg := fmt.Sprintf(" et - %s | %s to quit", fnameStr, quitKeyBindsString)
	for i, ch := range statusMsg {
		if i >= w {
			break
		}
		s.SetContent(i, h-1, ch, nil, statusStyle)
	}
	s.Show()
}

func IsKeyAny(key tcell.Key, keyAsRune string, keys []config.Key) bool {
	for _, k := range keys {
		if key == k.Key || (keyAsRune != "" && keyAsRune == k.String()) {
			return true
		}
	}
	return false
}

func main() {
	filename := flag.String("f", "", "file to open")
	showHelp := flag.Bool("help", false, "show help")
	showVersion := flag.Bool("version", false, "show version")
	showVersion2 := flag.Bool("v", false, "show version")
	flag.Parse()
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to create log file: %s", err.Error())
	}
	jsonSlog := slog.New(slog.NewJSONHandler(f, nil))
	slog.SetDefault(jsonSlog)

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	} else if *showVersion || *showVersion2 {
		fmt.Printf("et v%s\n", version)
		os.Exit(0)
	}
	if *filename == "" && len(os.Args) > 1 {
		*filename = os.Args[1]
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		slog.Error("creating screen", "err", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		slog.Error("initializing screen", "err", err)
		os.Exit(1)
	}
	defer screen.Fini()

	cfg := config.NewDefault()
	draw(screen, cfg, *filename)

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			draw(screen, cfg, *filename)
		case *tcell.EventKey:
			keyAsRune := ""
			if e.Key() == tcell.KeyRune {
				keyAsRune = e.Str()
			}
			if IsKeyAny(e.Key(), keyAsRune, cfg.KeyBindings.Quit) {
				return
			}
		}
	}
}
