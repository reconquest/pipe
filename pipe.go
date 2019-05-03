package pipe

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
)

var (
	debugPipeLinkedList = os.Getenv("DEBUG_PIPE_LINKED_LIST") == "1"
)

type Flower func(Flow)

func Pipe(flowers ...Flower) *Pipeline {
	pipeline := &Pipeline{
		mutex:    &sync.Mutex{},
		routines: &sync.WaitGroup{},
	}

	tail := &Flow{pipeline: pipeline}

	flows := []Flow{}
	for range flowers {
		flow := Flow{
			id:       int(atomic.AddInt64(&pipeline.flowID, 1)),
			streams:  map[string]*Stream{},
			pipeline: pipeline,
			process:  &sync.WaitGroup{},
		}

		// dereference the tail
		prev := *tail
		// set current element as next item of tail
		prev.next = &flow

		// set current tail as previous of current element
		flow.prev = &prev

		// set tail as reference to current element
		tail = &flow

		flow.process.Add(1)

		flows = append(flows, flow)
	}

	if debugPipeLinkedList {
		for tail != nil {
			log.Printf(
				"%d %p < %d %p > %d %p",
				tail.prev.ID(),
				tail.prev,
				tail.id,
				&tail,
				tail.next.ID(),
				tail.next,
			)

			tail = tail.prev
		}
	}

	pipeline.tail = flows[len(flows)-1]

	for i, flower := range flowers {
		flow := flows[i]

		pipeline.routines.Add(1)
		go func(flower Flower) {
			defer pipeline.routines.Done()
			defer flow.process.Done()

			flower(flow)

			flow.close()
		}(flower)
	}

	return pipeline
}

func EachLine(stream string, flower func(Flow, string)) Flower {
	return func(flow Flow) {
		for {
			line, ok := flow.Scan(stream)
			if !ok {
				return
			}

			flower(flow, line)
		}
	}
}

func Buffered(stream string, flower func(Flow, string)) Flower {
	return func(flow Flow) {
		contents, err := ioutil.ReadAll(flow.In(stream))
		if err != nil {
			log.Println(err)
			return
		}

		flower(flow, string(contents))
	}
}

func Exec(name string, args ...string) Flower {
	return func(flow Flow) {
		cmd := exec.Command(name, args...)

		cmd.Stdin = flow
		cmd.Stdout = flow

		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func Shell(cmdline string) Flower {
	return Exec(os.Getenv("SHELL"), "-c", cmdline)
}

func DumpIn() Flower {
	return EachLine(StreamStd, func(flow Flow, line string) {
		log.Println(line)
	})
}
