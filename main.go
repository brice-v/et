package main

import (
	"et/config"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
)

const logFileName = "et.log"
const version = "0.0.1"

func drawLine(s tcell.Screen, baseStyle tcell.Style, w, lineNumberOnScreen int, line string) {
	lineRunes := []rune(line)
	lineLen := len(lineRunes)
	for x := range w {
		ch := ' '
		if x < lineLen && lineLen != 0 {
			ch = lineRunes[x]
		}
		s.SetContent(x, lineNumberOnScreen, ch, nil, baseStyle)
	}
}

func drawContent(s tcell.Screen, baseStyle tcell.Style, w, h int, fileContent string) {
	lines := strings.Split(fileContent, "\n")
	numLines := len(lines)
	for i := range h {
		l := ""
		if i < numLines {
			l = lines[i]
		}
		drawLine(s, baseStyle, w, i, l)
	}
}

func drawStatusBar(s tcell.Screen, cfg *config.Config, baseStyle tcell.Style, w, h int, fileName string) {
	statusStyle := baseStyle.Background(cfg.Colors.StatusBar.Color)
	statusBarH := h - 1
	for x := range w {
		s.SetContent(x, statusBarH, ' ', nil, statusStyle)
	}
	quitKeyBindsString := cfg.GetQuitKeyBindingsAsStr()
	fnameStr := fileName
	if fileName == "" {
		fnameStr = "<new file>"
	}
	statusMsg := fmt.Sprintf(" et - %s | %s to quit", fnameStr, quitKeyBindsString)
	for i, ch := range statusMsg {
		if i >= w {
			break
		}
		s.SetContent(i, statusBarH, ch, nil, statusStyle)
	}
}

func draw(s tcell.Screen, cfg *config.Config, fileName, fileContent string) {
	style := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	s.Clear()
	w, h := s.Size()
	drawContent(s, style, w, h, fileContent)
	drawStatusBar(s, cfg, style, w, h, fileName)
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
	fileName := flag.String("f", "", "file to open")
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
	if *fileName == "" && len(os.Args) > 1 {
		*fileName = os.Args[1]
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
	fileContent := ""
	if *fileName != "" {
		data, err := os.ReadFile(*fileName)
		if err != nil {
			slog.Warn("could not read file", "err", err)
		}
		fileContent = string(data)
	}
	draw(screen, cfg, *fileName, fileContent)

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			draw(screen, cfg, *fileName, fileContent)
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
