package pipe

import "io"
import "io/ioutil"
import "sync"

type Flow struct {
	pipeline *Pipeline

	prev *Flow
	next *Flow

	id int

	streams map[string]*Stream

	io.Reader
	io.Writer

	process *sync.WaitGroup
}

func (flow *Flow) Wait() {
	flow.process.Wait()
}

func (flow Flow) Read(buffer []byte) (int, error) {
	return flow.In(StreamStd).Read(buffer)
}

func (flow Flow) Write(buffer []byte) (int, error) {
	return flow.Out(StreamStd).Write(buffer)
}

func (flow *Flow) In(name string) *Stream {
	return flow.pipeline.stream(flow, name, StreamModeRead)
}

func (flow *Flow) Out(name string) *Stream {
	return flow.pipeline.stream(flow, name, StreamModeWrite)
}

func (flow *Flow) ID() int {
	if flow == nil {
		return 0
	}

	return flow.id
}

func (flow *Flow) stream(name string) *Stream {
	stream, ok := flow.streams[name]
	if !ok {
		stream = newStream(name, flow.ID())
		flow.streams[name] = stream
	}

	return stream
}

func (flow *Flow) Scan(name string) (string, bool) {
	scanner := flow.In(name).Scanner()
	scanned := scanner.Scan()
	if scanned {
		text := scanner.Text()
		return text, scanned
	}

	return "", scanned
}

func (flow *Flow) close() {
	flow.pipeline.mutex.Lock()
	defer flow.pipeline.mutex.Unlock()

	for _, stream := range flow.streams {
		err := stream.Writer.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (flow *Flow) ReadAll(name string) string {
	contents, err := ioutil.ReadAll(flow.In(name))
	if err != nil {
		panic(err)
	}

	return string(contents)
}
