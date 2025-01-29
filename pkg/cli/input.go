package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func askForInput(prompt string) string {
	var input string

	fmt.Print(prompt)
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.TrimSpace(input)
	return input
}
