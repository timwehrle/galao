package galao

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type outgoingMessage struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Event string `json:"event"`
	Error string `json:"error,omitempty"`
}

type process struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser

	rd *bufio.Reader

	mu     sync.Mutex
	closed bool
}

func startProcess(ctx context.Context, binaryPath string) (*process, error) {
	cmd := exec.CommandContext(ctx, binaryPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &process{
		cmd:   cmd,
		stdin: stdin,
		rd:    bufio.NewReader(stdoutPipe),
	}, nil
}

func (p *process) send(v any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return io.ErrClosedPipe
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(payload)))

	if _, err := p.stdin.Write(lenBuf[:]); err != nil {
		return err
	}
	_, err = p.stdin.Write(payload)
	return err
}

func (p *process) readMessage(out any) error {
	var lenBuf [4]byte
	if _, err := io.ReadFull(p.rd, lenBuf[:]); err != nil {
		return err
	}

	n := binary.BigEndian.Uint32(lenBuf[:])
	if n == 0 {
		return fmt.Errorf("invalid frame length 0")
	}

	payload := make([]byte, n)
	if _, err := io.ReadFull(p.rd, payload); err != nil {
		return err
	}

	return json.Unmarshal(payload, out)
}

func (p *process) close() error {
	p.mu.Lock()
	p.closed = true
	defer p.mu.Unlock()

	_ = p.stdin.Close()
	return nil
}
