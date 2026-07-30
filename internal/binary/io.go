package binary

import (
	"io"
	"os"
)

// ReadFrom streams data from an io.Reader directly into the emitter's buffer.
func (e *Emitter) ReadFrom(r io.Reader) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	e.data = append(e.data, data...)
	return int64(len(data)), nil
}

// WriteTo implements io.WriterTo, flushing the emitted data to an io.Writer.
func (e *Emitter) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(e.data)
	return int64(n), err
}

// Write implements io.Writer, allowing the Emitter to be used as a streaming destination.
func (e *Emitter) Write(p []byte) (int, error) {
	e.Bytes(p)
	return len(p), nil
}

// WriteString implements io.StringWriter, appending string bytes directly.
func (e *Emitter) WriteString(s string) (int, error) {
	e.String(s)
	return len(s), nil
}

// WriteFile writes the complete binary buffer to disk at the specified path.
func (e *Emitter) WriteFile(path string) error {
	return os.WriteFile(path, e.data, 0o644)
}
