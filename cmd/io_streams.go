package cmd

import (
	"bytes"
	"io"
)

// IOStreams stores the standard names for iostreams.
type IOStreams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewTestIOStreams returns a valid IOStreams and in, out, errout buffers for testing.
func NewTestIOStreams() (IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return IOStreams{
		Stdin:  in,
		Stdout: out,
		Stderr: errOut,
	}, in, out, errOut
}
