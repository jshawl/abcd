# diffrn

## Vision

### CLI

```
diffn [--staged] [treeish]
```

opens a pager, watches for changes, and updates the diff

### Library

```go
package main

import (
    "fmt"
    "os"

    diffrn "github.com/jshawl/diffrn"
)

func main() {
    cmd := exec.Command("git", "diff", "--color")
	stdout, err := cmd.Output()
    diff := diffrn.Parse(string(stdout))

    fmt.Println(len(diff.Files))
    fmt.Println(len(diff.Files[0].Blocks))
    fmt.Println(len(diff.Files[0].Blocks[0].Lines))
}
```
