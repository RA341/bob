package logger

import (
	"bytes"
	"io"
)

type PrefixWriter struct {
	prefix []byte
	w      io.Writer
}

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{w: w, prefix: []byte(prefix)}
}

func (p *PrefixWriter) Write(b []byte) (n int, err error) {
	lines := bytes.Split(b, []byte("\n"))
	for i, line := range lines {
		// avoid adding prefix to the trailing empty element after a final newline
		if i == len(lines)-1 && len(line) == 0 {
			break
		}
		_, err = p.w.Write(append(p.prefix, line...))
		if err != nil {
			return
		}
		if i < len(lines)-1 {
			_, err = p.w.Write([]byte("\n"))
			if err != nil {
				return
			}
		}
	}
	return len(b), nil
}
