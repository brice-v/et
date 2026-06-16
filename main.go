package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"et/action"
	"et/defaults"

	"github.com/gdamore/tcell/v3"
)

func init() {
	if f, err := initLogging(); err != nil {
		log.Printf("warning: could not initialize log file: %s", err)
	} else {
		log.Printf("et %s started, logging to %s", version, f.Name())
	}
}

func main() {
	filename := flag.String("f", "", "file to open")
	showHelp := flag.Bool("help", false, "show help")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("et version %s\n", version)
		os.Exit(0)
	}

	if *filename == "" {
		if args := flag.Args(); len(args) > 0 {
			*filename = args[0]
		}
	}

	cfg := loadConfig()

	keyActions := make(map[keySpec]action.Action)
	for name, actStr := range cfg.Keybindings {
		ks, err := parseKeySpec(name)
		if err != nil {
			log.Fatalf("config error:\n  %s", err)
		}
		act, err := action.Parse(actStr)
		if err != nil {
			log.Fatalf("config error:\n  %s", err)
		}
		keyActions[ks] = act
	}

	var fileLines []string
	hasFile := *filename != ""
	if hasFile {
		data, err := os.ReadFile(*filename)
		if err != nil {
			log.Fatalf("error opening %s:\n  %s", *filename, err)
		}
		fileLines = strings.Split(string(data), "\n")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("screen error:\n  %s", err.Error())
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("screen init error:\n  %s", err.Error())
	}
	defer screen.Fini()

	bg := defaults.ColorBackground()
	if cfg.Colors.Background != "" {
		if c, err := parseColor(cfg.Colors.Background); err == nil {
			bg = c
		}
	}
	fg := defaults.ColorForeground()
	if cfg.Colors.Foreground != "" {
		if c, err := parseColor(cfg.Colors.Foreground); err == nil {
			fg = c
		}
	}
	status := defaults.ColorStatus()
	if cfg.Colors.StatusBG != "" {
		if c, err := parseColor(cfg.Colors.StatusBG); err == nil {
			status = c
		}
	}

	style := tcell.StyleDefault.Background(bg).Foreground(fg)
	statusStyle := style.Background(status)

	draw := func() {
		drawScreen(screen, *filename, hasFile, fileLines, style, statusStyle, cfg.TabWidth)
	}

	draw()

	for ev := range screen.EventQ() {
		switch e := ev.(type) {
		case *tcell.EventResize:
			draw()
		case *tcell.EventKey:
			for ks, act := range keyActions {
				if ks.matches(e) && act == action.Quit {
					return
				}
			}
		}
	}
}
