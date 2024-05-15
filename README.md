# diffrn

## CLI

```
diffn [--staged] [treeish]
```

opens a pager, watches for changes, and updates the diff

## Library

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    diffrn "github.com/jshawl/diffrn"
)

func main() {
    cmd := exec.Command("git", "diff")
    stdout, err := cmd.Output()
    diff := diffrn.ParseDiff(string(stdout))

    diffJson, _ := json.MarshalIndent(diff, "", "    ")
    fmt.Println(string(diffJson))
}
```
