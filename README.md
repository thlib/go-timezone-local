### Get the full name of the local timezone

Usage:

```go
package main

import (
    "fmt"
    "github.com/thlib/go-timezone-local/tzlocal"
)

func main() {
    val, _ := tzlocal.RuntimeTZ()
    fmt.Println(val)
}
```

All credit goes to [colm.anseo](https://stackoverflow.com/users/1218512/colm-anseo) and [MrFuppes](https://stackoverflow.com/users/10197418/mrfuppes) for providing the following answers:  
* https://stackoverflow.com/a/68938947/175071
* https://stackoverflow.com/a/68966317/175071

Installation
-----

```
go mod init example.com/yourpackage
go mod vendor
go run main.go
```
