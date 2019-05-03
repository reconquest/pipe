package pipe

import (
	"bufio"
	"io"
	"os"
	"sync"
)

type (
	StreamMode byte
)

const (
	StreamModeRead  StreamMode = 0xA
	StreamModeWrite StreamMode = 0xF

	StreamStd = "std"
)

func (mode StreamMode) String() string {
	switch mode {
	case StreamModeRead:
		return "read"
	case StreamModeWrite:
		return "write"
	}
	return "unknown"
}

type Stream struct {
	*sync.Mutex
	flowID  int
	Reader  io.ReadCloser
	Writer  io.WriteCloser
	scanner *bufio.Scanner

	name string
}

func (stream *Stream) Read(buffer []byte) (int, error) {
	return stream.Reader.Read(buffer)
}

func (stream *Stream) Write(buffer []byte) (int, error) {
	return stream.Writer.Write(buffer)
}

func (stream *Stream) Scanner() *bufio.Scanner {
	stream.Lock()
	defer stream.Unlock()

	if stream.scanner == nil {
		stream.scanner = bufio.NewScanner(stream.Reader)
	}

	return stream.scanner
}

func newStream(name string, flowID int) *Stream {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stream := &Stream{
		name:   name,
		flowID: flowID,
		Reader: reader,
		Mutex:  &sync.Mutex{},
		Writer: writer,
	}

	return stream
}
