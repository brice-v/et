//go:build !windows

package terminal

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
)

type unixPtyBackend struct {
	cmd *exec.Cmd
	f   *os.File
}

func newPtyBackend(cmd *exec.Cmd, cols, rows uint16) (ptyBackend, error) {
	winsize := pty.Winsize{
		Cols: cols,
		Rows: rows,
	}
	f, err := pty.StartWithAttrs(
		cmd,
		&winsize,
		&syscall.SysProcAttr{
			Setsid:  true,
			Setctty: true,
			Ctty:    1,
		})
	if err != nil {
		return nil, err
	}
	return &unixPtyBackend{cmd: cmd, f: f}, nil
}

func (b *unixPtyBackend) Read(p []byte) (int, error)        { return b.f.Read(p) }
func (b *unixPtyBackend) Write(p []byte) (int, error)       { return b.f.Write(p) }
func (b *unixPtyBackend) WriteString(s string) (int, error) { return b.f.WriteString(s) }

func (b *unixPtyBackend) Close() error {
	if b.cmd != nil && b.cmd.Process != nil {
		if err := b.cmd.Process.Kill(); err != nil {
			slog.Warn("error killing command", "err", err)
		}
		if err := b.cmd.Wait(); err != nil {
			slog.Warn("error waiting for command", "err", err)
		}
	}
	return b.f.Close()
}

func (b *unixPtyBackend) Resize(cols, rows uint16) error {
	return pty.Setsize(b.f, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
}
