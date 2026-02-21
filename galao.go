package galao

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

//go:embed renderer/GalaoRenderer.app/**
var rendererBundle embed.FS

type EventHandler func(Event)

type Event struct {
	ID   string
	Name string
}

type App struct {
	proc     *process
	handlers map[string]EventHandler

	cancel context.CancelFunc
	waitCh chan error
}

func New() *App {
	return &App{handlers: make(map[string]EventHandler)}
}

func (a *App) SetView(node ViewNode) error {
	if a.proc == nil {
		return errors.New("renderer not started")
	}
	return a.proc.send(struct {
		Type string   `json:"type"`
		Tree ViewNode `json:"tree"`
	}{
		Type: "set_view",
		Tree: node,
	})
}

func (a *App) OnEvent(id string, handler EventHandler) {
	a.handlers[id] = handler
}

func (a *App) Close() error {
	if a.cancel != nil {
		a.cancel()
	}
	if a.proc != nil {
		_ = a.proc.close()
	}
	if a.waitCh != nil {
		<-a.waitCh
	}
	return nil
}

func (a *App) Run(ctx context.Context, setup func()) error {
	tmp, err := os.MkdirTemp("", "galao-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := extractBundle(tmp); err != nil {
		return err
	}

	macOSDir := filepath.Join(tmp, "GalaoRenderer.app", "Contents", "MacOS")
	entries, err := os.ReadDir(macOSDir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("no executable found in %s", macOSDir)
	}
	exe := filepath.Join(macOSDir, entries[0].Name())
	_ = os.Chmod(exe, 0755)

	runCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	a.proc, err = startProcess(runCtx, exe)
	if err != nil {
		return err
	}

	a.waitCh = make(chan error, 1)
	go func() {
		a.waitCh <- a.proc.cmd.Wait()
	}()

	if err := a.waitReady(runCtx, 5*time.Second); err != nil {
		_ = a.Close()
		return err
	}

	setup()

	for {
		select {
		case <-runCtx.Done():
			_ = a.Close()
			return runCtx.Err()
		case err := <-a.waitCh:
			if err == nil {
				return io.EOF
			}
			return err
		default:
			var msg outgoingMessage
			if err := a.proc.readMessage(&msg); err != nil {
				// If process exited concurrently, surface that error instead of the read error
				select {
				case werr := <-a.waitCh:
					if werr != nil {
						return werr
					}
					return err
				default:
				}
				return err
			}

			switch msg.Type {
			case "event":
				if h, ok := a.handlers[msg.ID]; ok {
					h(Event{ID: msg.ID, Name: msg.Event})
				}
			case "error":
				return fmt.Errorf("renderer error: %s", msg.Error)
			case "ready":
				// Ignore, already handled in waitReady
			}
		}
	}
}

func (a *App) waitReady(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("renderer not ready: %w", ctx.Err())
		case err := <-a.waitCh:
			if err == nil {
				return fmt.Errorf("renderer exited before ready")
			}
			return fmt.Errorf("renderer exited before ready: %w", err)
		default:
			var msg outgoingMessage
			if err := a.proc.readMessage(&msg); err != nil {
				return err
			}
			if msg.Type == "ready" {
				return nil
			}
			if msg.Type == "error" {
				return fmt.Errorf("renderer error: %s", msg.Error)
			}
		}
	}
}

func extractBundle(dst string) error {
	return fs.WalkDir(rendererBundle, "renderer/GalaoRenderer.app", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("renderer", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := rendererBundle.Open(path)
		if err != nil {
			return err
		}
		defer data.Close()

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		out, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, data); err != nil {
			return err
		}

		if filepath.Base(filepath.Dir(targetPath)) == "MacOS" {
			return os.Chmod(targetPath, 0755)
		}
		return nil
	})
}
