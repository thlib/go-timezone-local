### Get the full name of the local timezone

```
go get github.com/thlib/go-timezone-local/tzlocal
```

See it in action:
-----

Open your project folder  
Create a file `main.go`

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

Run the following commands:
```
go mod init example.com/yourpackage
go mod vendor
go run main.go
```

It should print your OS timezone.

For developers of github.com/thlib/go-timezone-local/tzlocal, updating the list of time zones in windows
-----

Clone github.com/thlib/go-timezone-local  
Change directory to go-timezone-local  

```
cd go-timezone-local
go generate ./...
```

Credit
------

All credit goes to [colm.anseo](https://stackoverflow.com/users/1218512/colm-anseo) and [MrFuppes](https://stackoverflow.com/users/10197418/mrfuppes) for providing the following answers:  
* https://stackoverflow.com/a/68938947/175071
* https://stackoverflow.com/a/68966317/175071
