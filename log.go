package main

import (
	"log"
	"os"
	"path/filepath"
)

var cfgLog *log.Logger

func initLogging() (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	etDir := filepath.Join(configDir, "et")
	if err := os.MkdirAll(etDir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(etDir, "et.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	cfgLog = log.New(f, "config: ", log.Ltime)
	return f, nil
}

func warn(format string, args ...any) {
	if cfgLog != nil {
		cfgLog.Printf("warning: "+format, args...)
	}
}
