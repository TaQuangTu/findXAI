package liberror

import (
	"fmt"
)

type AdditionalData map[string]any

func (d AdditionalData) Merge(data AdditionalData) AdditionalData {
	newData := make(map[string]any)
	for key, vData := range d {
		newData[key] = vData
	}
	for key, vData := range data {
		newData[key] = vData
	}
	return newData
}

func (a AdditionalData) String() string {
	dataStr := ""
	for key, value := range a {
		dataStr += fmt.Sprintf("%s: %v\n", key, value)
	}
	return dataStr
}

type Severity string
