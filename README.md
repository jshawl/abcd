# abcd

A Better Console Diff

## CLI

```
abcd [--staged] [treeish]
```

opens a pager, watches for changes, and updates the diff

## Library

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    abcd "github.com/jshawl/abcd"
)

func main() {
    cmd := exec.Command("git", "diff")
    stdout, err := cmd.Output()
    diff := abcd.ParseDiff(string(stdout))

    diffJson, _ := json.MarshalIndent(diff, "", "    ")
    fmt.Println(string(diffJson))
}
```
