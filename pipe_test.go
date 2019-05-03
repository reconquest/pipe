package pipe

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipeFunctions(t *testing.T) {
	test := assert.New(t)

	_ = test

	lister := Exec("ls", "/etc/")

	filter := EachLine(
		StreamStd,
		func(flow Flow, line string) {
			if strings.Contains(line, ".") {
				return
			}

			if len(line) > 3 {
				return
			}

			fmt.Fprintln(flow, line)
		},
	)

	sorter := Buffered(
		StreamStd,
		func(flow Flow, buffer string) {
			lines := strings.Split(string(buffer), "\n")

			sort.Strings(lines)

			for i := 0; i < len(lines); i++ {
				fmt.Println(lines[i])
			}
		},
	)

	Pipe(
		lister,
		filter,
		sorter,
	).Wait()
}

func TestPipeCommands(t *testing.T) {
	test := assert.New(t)

	_ = test

	Pipe(
		Exec("ls", "/etc"),
		Exec("cut", "-d.", "-f2"),
		Exec("sort"),
		Exec("uniq", "-c"),
		Exec("sort", "-n"),
		Exec("tail", "-n10"),
		DumpIn(),
	).Wait()
}

func TestCustomStreams(t *testing.T) {
	test := assert.New(t)

	_ = test

	Pipe(
		func(flow Flow) {
			flow.Out("email").Write([]byte("data in email"))
		},

		func(flow Flow) {
			data, err := ioutil.ReadAll(flow.In("email"))
			if err != nil {
				panic(err)
			}

			log.Printf("data from 'email' stream: %q", string(data))
		},
	).Wait()
}

func TestNoWait(t *testing.T) {
	test := assert.New(t)

	_ = test

	calls := []int{}

	pipe := Pipe(
		func(flow Flow) {
			time.Sleep(time.Millisecond * 100)

			calls = append(calls, flow.id)

			flow.Out("aaa").Write([]byte("aaa data"))
		},
		func(flow Flow) {
			time.Sleep(time.Millisecond * 50)

			calls = append(calls, flow.id)

			buffer := flow.ReadAll("aaa")

			flow.Out("bbb").Write([]byte("bbb data: " + buffer))
		},
	)

	calls = append(calls, 0)

	out := pipe.Out("bbb")

	test.EqualValues([]int{0, 2, 1}, calls)

	test.EqualValues("bbb data: aaa data", out)
}
