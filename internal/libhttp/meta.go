package libhttp

import (
	"io"
	"time"
)

type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
)

type EncoderType int

const (
	JSON EncoderType = iota
)

type (
	EncoderOption struct {
		Type EncoderType
		Data any
	}

	RequestOption struct {
		RequestTimeout time.Duration
		Method         Method
		Url            string
		Header         Header
		Body           *EncoderOption
	}

	Encoder interface {
		Encode(data any) (io.Reader, error)
	}
)
