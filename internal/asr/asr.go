package asr

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
)

type ASR struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	available bool
	errorCh   chan error // Канал для получения ошибок из stderr

}

func NewASR() *ASR {
	cmd := exec.Command("path/to/asr") // Заменить на команду для запуска приложения ASR

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	asr := &ASR{
		cmd:       cmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		available: true,
		errorCh:   make(chan error), // Initialize the error channel
	}

	return asr
}

func (a *ASR) IsAvailable() bool {
	return a.available
}

func (a *ASR) Start() error {
	err := a.cmd.Start()
	if err != nil {
		return err
	}

	a.startErrorListener()
	// Ожидаем ошибку
	select {
	case err := <-a.errorCh:
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ASR) Stop() error {
	err := a.stdin.Close()
	if err != nil {
		return err
	}

	err = a.cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (a *ASR) ProcessAudioChunk(chunk []byte) (string, error) {
	_, err := a.stdin.Write(chunk)
	if err != nil {
		return "", err
	}

	responseBytes, err := a.readASRResponse()
	if err != nil {
		return "", err
	}

	var response struct {
		Text string `json:"text"`
	}

	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return "", err
	}

	return response.Text, nil
}

func (a *ASR) readASRResponse() ([]byte, error) {
	responseBytes := make([]byte, 0)
	buffer := make([]byte, 4096)

	for {
		n, err := a.stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		responseBytes = append(responseBytes, buffer[:n]...)
	}

	// Пример responseBytes = []byte(`{"text": "Sample ASR response"}`)
	return responseBytes, nil
}

func (a *ASR) startErrorListener() {
	go func() {
		scanner := bufio.NewScanner(a.stderr)
		for scanner.Scan() {
			errMsg := scanner.Text()
			a.errorCh <- errors.New(errMsg)
		}
	}()
}

// GetErrorCh returns the channel for receiving errors from stderr
func (a *ASR) GetErrorCh() <-chan error {
	return a.errorCh
}
