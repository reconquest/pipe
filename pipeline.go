package pipe

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/reconquest/nopio-go"
)

var (
	debugPipelineStream = os.Getenv("DEBUG_PIPELINE_STREAM") == "1"
)

type pipeline struct {
	mutex *sync.Mutex
}

var (
	nopStream = &Stream{
		name:   "closed",
		id:     -1,
		Reader: nopio.NopReadCloser{},
		Writer: nopio.NopWriteCloser{},
	}
)

func (pipeline *pipeline) stream(flow *Flow, name string, mode StreamMode) *Stream {
	pipeline.mutex.Lock()
	defer pipeline.mutex.Unlock()

	var stream *Stream
	var ok bool
	switch mode {
	case StreamModeRead:
		if flow.prev.streams == nil {
			return nopStream
		}

		stream, ok = flow.prev.streams[name]
		if !ok {
			stream = newStream(name)
			flow.prev.streams[name] = stream
		}

	case StreamModeWrite:
		stream, ok = flow.streams[name]
		if !ok {
			stream = newStream(name)
			flow.streams[name] = stream
		}

	default:
		panic(fmt.Errorf("unexpected stream mode: %v", mode))
	}

	if debugPipelineStream {
		log.Printf(
			"flow.id: %v, stream: %d %p, mode: %s",
			flow.id,
			stream.id,
			stream,
			mode,
		)
	}

	return stream
}
