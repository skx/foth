// part3 driver

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	forth := NewEval()

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

		forth.Eval(strings.Split(text, " "))
	}
}
