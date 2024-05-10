package main

import (
	"fmt"
	"os/exec"
)

type Diff struct {
}

func main() {
	cmd := exec.Command("git", "diff", "--color")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}
