package tmux

import (
	"context"
	"fmt"
	"strings"
)

// SendEnter sends an Enter keystroke to a tmux session.
func (m *Manager) SendEnter(ctx context.Context, name string) error {
	return m.SendKeys(ctx, name, "Enter")
}

// SendKeys sends arbitrary key sequences to a tmux session.
func (m *Manager) SendKeys(ctx context.Context, name string, keys ...string) error {
	sessionName := SanitizeName(name)
	args := append([]string{"send-keys", "-t", sessionName}, keys...)
	res, err := m.exec.Run(ctx, "tmux", args...)
	if err != nil {
		return fmt.Errorf("tmux send-keys: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("tmux send-keys failed (exit %d): %s", res.ExitCode, res.Stderr)
	}
	return nil
}

// SendText sends text to a tmux session followed by Enter.
// For multi-line text, each line is sent separately with Enter after each.
func (m *Manager) SendText(ctx context.Context, name string, text string) error {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if err := m.SendKeys(ctx, name, line, "Enter"); err != nil {
			return err
		}
	}
	return nil
}

// SendCommand sends a command string to a tmux session (typed then Enter).
func (m *Manager) SendCommand(ctx context.Context, name string, cmd string) error {
	return m.SendKeys(ctx, name, cmd, "Enter")
}
