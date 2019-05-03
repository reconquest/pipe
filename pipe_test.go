package pipe

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoNameYet(t *testing.T) {
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
		func(flow Flow, buffer []byte) {
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
	)
}

func TestNoNameYet2(t *testing.T) {
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
	)
}
