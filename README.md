### Get the full name of the local timezone

Usage:

```go
package main

import (
    "fmt"
    tz "github.com/thlib/go-local-timezone"
)

func main() {
    val, _ := tz.RuntimeTZ()
    fmt.Println(val)
}
```
