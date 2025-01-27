package cli

import (
	"bufio"
	"fmt"
	"os"
)

func askForInput(prompt string) string {
	var input string

	fmt.Print(prompt)
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	return input
}
