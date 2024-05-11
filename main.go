package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("git", "diff")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	diff, _ := parseDiff(string(stdout))
	diffJson, _ := json.MarshalIndent(diff, "", "    ")
	fmt.Println(string(diffJson))
}
