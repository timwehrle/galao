package galao

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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
}

func New() *App {
	return &App{handlers: make(map[string]EventHandler)}
}

func (a *App) SetView(node ViewNode) error {
	return a.proc.send(map[string]any{
		"type": "set_view",
		"tree": node,
	})
}

func (a *App) OnEvent(id string, handler EventHandler) {
	a.handlers[id] = handler
}

func (a *App) Run(setup func()) error {
	tmp, err := os.MkdirTemp("", "galao-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := extractBundle(tmp); err != nil {
		return err
	}

	binaryPath := filepath.Join(tmp, "GalaoRenderer.app", "Contents", "MacOS", "GalaoRenderer")

	a.proc, err = startProcess(binaryPath)
	if err != nil {
		return err
	}

	setup()

	for {
		msg, err := a.proc.readLine()
		if err != nil {
			return err
		}
		if msg.Type == "event" {
			if h, ok := a.handlers[msg.ID]; ok {
				h(Event{ID: msg.ID, Name: msg.Event})
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

		if filepath.Base(targetPath) == "GalaoRenderer" {
			return os.Chmod(targetPath, 0755)
		}

		return nil
	})
}
