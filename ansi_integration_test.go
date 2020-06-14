package ansi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aoldershaw/ansi"
	. "github.com/onsi/gomega"
)

func TestAnsi_Integration_InMemory(t *testing.T) {
	for _, tt := range []struct {
		description string
		events      [][]byte
		lines       []ansi.Line
	}{
		{
			description: "basic test",
			events: [][]byte{
				[]byte("hello\nworld"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("hello"),
					},
				},
				{
					{
						Data: ansi.Text("world"),
					},
				},
			},
		},
		{
			description: "styling",
			events: [][]byte{
				[]byte("hello \x1b[1mworld\x1b[m\n"),
				[]byte("\x1b[31mthis is red\x1b[m\n"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("hello "),
					},
					{
						Data:  []byte("world"),
						Style: ansi.Style{Bold: true},
					},
				},
				{
					{
						Data:  []byte("this is red"),
						Style: ansi.Style{Foreground: ansi.Red},
					},
				},
			},
		},
		{
			description: "control sequences split over multiple events",
			events: [][]byte{
				[]byte("\x1b[31mthis is red\x1b"),
				[]byte("[0m but this is not"),
			},
			lines: []ansi.Line{
				{
					{
						Data:  []byte("this is red"),
						Style: ansi.Style{Foreground: ansi.Red},
					},
					{
						Data: ansi.Text(" but this is not"),
					},
				},
			},
		},
		{
			description: "moving the cursor",
			events: [][]byte{
				[]byte("hello\x1b[3Cworld"),
				[]byte("\x1b[Ggoodbye"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("goodbye world"),
					},
				},
			},
		},
		{
			description: "save and restore cursor",
			events: [][]byte{
				[]byte("\x1b[shello   world"),
				[]byte("\x1b[ugoodbye"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("goodbye world"),
					},
				},
			},
		},
		{
			description: "erase line",
			events: [][]byte{
				[]byte("this text is very important and will never be removed!\n"),
				[]byte("\x1b[1A\x1b[2Knevermind"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("nevermind"),
					},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)

			out := &ansi.InMemory{}
			log := ansi.New(out)

			initialEvents := make([][]byte, len(tt.events))
			for i, evt := range tt.events {
				initialEvents[i] = make([]byte, len(evt))
				copy(initialEvents[i], evt)
			}

			for _, evt := range tt.events {
				log.Parse(evt)
			}

			g.Expect(out.Lines).To(Equal(tt.lines))
			g.Expect(tt.events).To(Equal(initialEvents), "modified input bytes")
		})
	}
}

func Example() {
	output := &ansi.InMemory{}
	interpreter := ansi.New(output)

	interpreter.Parse([]byte("\x1b[1mbold\x1b[m text"))
	interpreter.Parse([]byte("\nline 2"))

	linesJSON, _ := json.Marshal(output.Lines)
	fmt.Println(string(linesJSON))
	// Output: [[{"data":"bold","style":{"bold":true}},{"data":" text","style":{}}],[{"data":"line 2","style":{}}]]
}
