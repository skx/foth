// part6 driver
//
// Loads "foth.4th" from cwd, if present, and evaluates it before the REPL
// is launched - otherwise the same as previous versions.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"foth/eval"
)

// If the given file exists, read the contents, and evaluate it
func doInit(eval *eval.Eval, path string) {

	handle, err := os.Open(path)
	if err != nil {
		return
	}

	reader := bufio.NewReader(handle)
	line, err := reader.ReadString(byte('\n'))
	for err == nil {

		// Trim it
		line = strings.TrimSpace(line)

		// Is this isn't comment then execute it
		if !strings.HasPrefix(line, "#") {

			// Evaluate
			err := eval.Eval(strings.Split(line, " "))
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			}
		}

		// Repeat
		line, err = reader.ReadString(byte('\n'))
	}

	handle.Close()
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	forth := eval.New()

	// Load the init-file if it is present.
	doInit(forth, "foth.4th")

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

		err = forth.Eval(strings.Split(text, " "))
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}

	}
}
