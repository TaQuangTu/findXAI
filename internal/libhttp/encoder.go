package libhttp

import (
	"bytes"
	"encoding/json"
	"io"
)

type JsonEncoder struct{}

func (e *JsonEncoder) Encode(data any) (io.Reader, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

var encoder = map[EncoderType]Encoder{
	JSON: &JsonEncoder{},
}
