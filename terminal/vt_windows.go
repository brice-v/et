//go:build windows

package terminal

import (
	"errors"
	"io"
	"log/slog"
	"os/exec"
)

type windowsPtyBackend struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func newPtyBackend(cmd *exec.Cmd, cols, rows uint16) (ptyBackend, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, err
	}

	return &windowsPtyBackend{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (b *windowsPtyBackend) Read(p []byte) (int, error) { return b.stdout.Read(p) }
func (b *windowsPtyBackend) Write(p []byte) (int, error) {
	return b.stdin.Write(p)
}
func (b *windowsPtyBackend) WriteString(s string) (int, error) {
	return b.stdin.Write([]byte(s))
}

func (b *windowsPtyBackend) Close() error {
	var firstErr error
	if b.cmd != nil && b.cmd.Process != nil {
		if err := b.cmd.Process.Kill(); err != nil {
			slog.Warn("error killing command", "err", err)
			if firstErr == nil {
				firstErr = err
			}
		}
		if err := b.cmd.Wait(); err != nil {
			slog.Warn("error waiting for command", "err", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	if err := b.stdin.Close(); err != nil {
		if firstErr == nil {
			firstErr = err
		}
	}
	if err := b.stdout.Close(); err != nil {
		if firstErr == nil {
			firstErr = err
		}
	}
	if b.stderr != nil {
		if err := b.stderr.Close(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (b *windowsPtyBackend) Resize(cols, rows uint16) error {
	return errors.New("resize not supported on Windows pipe backend")
}
