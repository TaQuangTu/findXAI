package libcmd

import (
	"context"
	"io"
)

type ICustomCmd interface {
	Run(ctx context.Context, args ...string) (stdoutOutput string, err error)
	WithStreamReader(stdoutStreamReader func(io.ReadCloser) error) ICustomCmd
}

type CustomCmd struct {
	Cmd *Cmd

	useStream func(io.ReadCloser) error
}

func (c *CustomCmd) Run(ctx context.Context, args ...string) (stdoutOutput string, err error) {
	c.Cmd.Args = args
	if c.useStream != nil {
		err = c.Cmd.RunStream(ctx, c.useStream)
		return
	}
	stdoutOutput, err = c.Cmd.Run(ctx, nil)
	return
}

func (c *CustomCmd) WithStreamReader(stdoutStreamReader func(io.ReadCloser) error) ICustomCmd {
	c.useStream = stdoutStreamReader
	return c
}

func NewCustomCmd(cmd string) ICustomCmd {
	return &CustomCmd{
		Cmd: &Cmd{
			Command: cmd,
		},
	}
}
