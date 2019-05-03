package pipe

import (
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
)

var (
	debugPipeLinkedList = false
)

type Flower func(Flow)

func Pipe(flowers ...Flower) {
	routines := &sync.WaitGroup{}

	pipeline := &pipeline{
		mutex: &sync.Mutex{},
	}

	tail := &Flow{pipeline: pipeline}

	flows := []Flow{}
	for _ = range flowers {
		flow := Flow{
			id:       getFlowID(),
			streams:  map[string]*Stream{},
			pipeline: pipeline,
		}

		// dereference the tail
		prev := *tail
		// set current element as next item of tail
		prev.next = &flow

		// set current tail as previous of current element
		flow.prev = &prev

		// set tail as reference to current element
		tail = &flow

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

	for i, flower := range flowers {
		flow := flows[i]

		routines.Add(1)
		go func(flower Flower) {
			defer routines.Done()
			flower(flow)

			flow.close()
		}(flower)
	}

	routines.Wait()
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

func Buffered(stream string, flower func(Flow, []byte)) Flower {
	return func(flow Flow) {
		contents, err := ioutil.ReadAll(flow)
		if err != nil {
			log.Println(err)
			return
		}

		flower(flow, contents)
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

func DumpIn() Flower {
	return EachLine(StreamStd, func(flow Flow, line string) {
		log.Println(line)
	})
}
