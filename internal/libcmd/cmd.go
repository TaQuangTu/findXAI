package libcmd

import (
	"bufio"
	"context"
	"findx/internal/liberror"
	"fmt"
	"io"
	"os/exec"
)

type Cmd struct {
	Command string
	Args    []string
}

type StdoutLineTextGetter func(string) error
type StdoutStreamReader func(io.ReadCloser) error

func (c *Cmd) Run(ctx context.Context, lineTextGetter StdoutLineTextGetter) (ouput string, err error) {
	cmd := exec.CommandContext(ctx, c.Command, c.Args...)
	stdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to get stdout pipe")
		return
	}
	stdErrReader, err := cmd.StderrPipe()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to get stderr pipe")
		return
	}
	err = cmd.Start()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to start command")
		return
	}
	var (
		stdoutScanner = bufio.NewScanner(stdOutReader)
		stderrScanner = bufio.NewScanner(stdErrReader)

		stdoutChan = make(chan string)
		stderrChan = make(chan string)
	)

	go func() {
		for stdoutScanner.Scan() {
			stdoutChan <- stdoutScanner.Text()
		}
		close(stdoutChan)
	}()

	go func() {
		for stderrScanner.Scan() {
			stderrChan <- stderrScanner.Text()
		}
		close(stderrChan)
	}()

	var (
		hasError = false
		errStr   string

		interruptedIssue = false
		interruptedError error
	)
	for {
		select {
		case lineText, ok := <-stdoutChan:
			if ok {
				ouput += fmt.Sprintf("%s\n", lineText)
				if !interruptedIssue && lineTextGetter != nil {
					interruptedError = lineTextGetter(lineText)
					if interruptedError != nil {
						interruptedIssue = true
					}
				}
			}
		case errOut, ok := <-stderrChan:
			if ok {
				errStr += fmt.Sprintf("%s\n", errOut)
				hasError = true
			}
		}
		if len(stdoutChan) == 0 && len(stderrChan) == 0 {
			break
		}
	}

	err = cmd.Wait()
	if err != nil {
		err = liberror.WrapMessage(err, "command execution failed")
		return
	}
	if hasError {
		err = liberror.WrapMessage(err, errStr)
		return
	}
	if interruptedIssue {
		err = liberror.WrapMessage(liberror.
			ErrorInterrupted.
			Wrap(interruptedError), "interrupted by getter")
		return
	}
	return
}

// RunStream speedup cmd task
// TODO: support error handling later
func (c *Cmd) RunStream(ctx context.Context, reader StdoutStreamReader) (err error) {
	if reader == nil {
		err = liberror.WrapMessage(liberror.ErrorDataInvalid, "stream getter is required")
		return
	}
	cmd := exec.CommandContext(ctx, c.Command, c.Args...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to get stdout pipe")
		return
	}
	err = cmd.Start()
	if err != nil {
		err = liberror.WrapMessage(err, "failed to start command")
		return
	}
	err = reader(stdoutReader)
	if err != nil {
		err = liberror.WrapMessage(err, "operate stdout stream failed")
		return
	}
	err = cmd.Wait()
	if err != nil {
		err = liberror.WrapMessage(err, "command execution failed")
		return
	}
	return
}
