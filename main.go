package main

import (
	"et/config"
	"et/keys"
	"et/ui"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v3"
	_ "github.com/gdamore/tcell/v3/encoding"
)

const logFileName = "et.log"
const version = "0.0.1"

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
	ui.Draw(screen, cfg, *fileName, fileContent)

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			ui.Draw(screen, cfg, *fileName, fileContent)
		case *tcell.EventKey:
			keyAsRune := ""
			if e.Key() == tcell.KeyRune {
				keyAsRune = e.Str()
			}
			if keys.IsKeyAny(e.Key(), keyAsRune, e.Modifiers(), cfg.KeyBindings.Quit) {
				return
			}
		}
	}
}
