package galao

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os/exec"
)

/* type incomingMessage struct {
	Type string    `json:"type"`
	Tree *ViewNode `json:"tree,omitempty"`
} */

type outgoingMessage struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Event string `json:"event"`
}

type process struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
}

func startProcess(binaryPath string) (*process, error) {
	ctx := context.Background()
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
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdoutPipe),
	}, nil
}

func (p *process) send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = p.stdin.Write(data)
	return err
}

func (p *process) readLine() (outgoingMessage, error) {
	var msg outgoingMessage
	if p.stdout.Scan() {
		err := json.Unmarshal(p.stdout.Bytes(), &msg)
		return msg, err
	}
	return msg, io.EOF
}
