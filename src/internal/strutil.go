package internal

import (
	"bytes"
	"fmt"
)

func MapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", k, v)
	}
	return b.String()
}
