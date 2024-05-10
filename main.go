package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

func main() {
	cmd := exec.Command("git", "diff", "--color")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	re := regexp.MustCompile("^diff")
	output := re.Split(string(stdout), -1)

	// Print the output
	fmt.Println(output)
}
