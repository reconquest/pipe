package pipe

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/reconquest/nopio-go"
)

var (
	debugPipelineStream = os.Getenv("DEBUG_PIPELINE_STREAM") == "1"
)

type Pipeline struct {
	mutex *sync.Mutex

	tail     Flow
	routines *sync.WaitGroup

	flowID int64
}

func (pipeline *Pipeline) Wait() *Pipeline {
	pipeline.routines.Wait()
	return pipeline
}

var (
	nopStream = &Stream{
		name:   "closed",
		flowID: 0,
		Reader: nopio.NopReadCloser{},
		Writer: nopio.NopWriteCloser{},
	}
)

func (pipeline *Pipeline) stream(flow *Flow, name string, mode StreamMode) *Stream {
	pipeline.mutex.Lock()
	defer pipeline.mutex.Unlock()

	var stream *Stream
	switch mode {
	case StreamModeRead:
		if flow.prev.streams == nil {
			return nopStream
		}

		stream = flow.prev.stream(name)

	case StreamModeWrite:
		stream = flow.stream(name)

	default:
		panic(fmt.Errorf("unexpected stream mode: %v", mode))
	}

	if debugPipelineStream {
		log.Printf(
			"flow.id: %v, stream: %d %p, mode: %s",
			flow.id,
			stream.flowID,
			stream,
			mode,
		)
	}

	return stream
}

func (pipeline *Pipeline) Stream(name string) *Stream {
	return pipeline.tail.stream(name)
}

func (pipeline *Pipeline) Out(name string) string {
	contents, err := ioutil.ReadAll(pipeline.Stream(name))
	if err != nil {
		panic(err)
	}

	return string(contents)
}
