// foth - final revision, allow if/else/then, neaten-code, and run files
//        specified on the command-line.  If none run the REPL.
//
// Loads "foth.4th" from cwd, if present, and evaluates it before the REPL
// is launched - otherwise the same as previous versions.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/skx/foth/foth/eval"
)

// "secret" word
func secret() error {
	fmt.Printf("nothing happens\n")
	return nil
}

// If the given file exists, read the contents, and evaluate it
func doInit(eval *eval.Eval, path string) error {

	handle, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(handle)
	line, err := reader.ReadString(byte('\n'))
	for err == nil {

		// Trim it
		line = strings.TrimSpace(line)

		// Evaluate
		err = eval.Eval(line)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())

			// Reset our state, to allow recovery
			eval.Reset()
		}

		// Repeat
		line, err = reader.ReadString(byte('\n'))
	}

	if err != nil {
		if err != io.EOF {
			return err
		}
	}

	err = handle.Close()
	if err != nil {
		return err
	}

	return nil
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	forth := eval.New()

	forth.Dictionary = append(forth.Dictionary, eval.Word{Name: "xyzzy", Function: secret})

	// Load the init-file if it is present.
	//
	// i.e. Run the file, but ignore errors.
	doInit(forth, "foth.4th")

	// If we got any arguments treat them as files to lead
	if len(os.Args) > 1 {
		for _, file := range os.Args[1:] {
			err := doInit(forth, file)
			if err != nil {
				fmt.Printf("error running %s: %s\n", file, err.Error())
				return
			}
		}
		return
	}

	// No arguments, just run the REPL
	for {
		fmt.Printf("> ")

		// Read input
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading input: %s\n", err.Error())
			return
		}

		// Trim it
		text = strings.TrimSpace(text)

		err = forth.Eval(text)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())

			// Reset our state, to allow recovery
			forth.Reset()
		}

	}
}
