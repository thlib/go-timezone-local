## Get the full name of your local time zone (OS setting)

### Why?
built-in functionality sometimes won't suffice:
```go
    zone, _ := time.Now().Zone() // try to get my time zone...
    loc, _ := time.LoadLocation(zone)
    fmt.Println(zone, loc) // prints e.g. CEST UTC -> obviously wrong!
```
localizing a date with obtained `loc` will cause
> panic: time: missing Location in call to Date

---

### Package Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/thlib/go-timezone-local/tzlocal"
)

func main() {
    tzname, err := tzlocal.RuntimeTZ()
    fmt.Println(tzname, err)

    // example:
    // tzname = "Europe/Berlin"

    // now you can use tzname to properly set up a location:
    loc, _ := time.LoadLocation(tzname)

    d0 := time.Date(2021, 10, 30, 20, 0, 0, 0, loc) // DST active:
    fmt.Println(d0)
    // 2021-10-30 20:00:00 +0200 CEST

    d1 := d0.AddDate(0, 0, 1) // add one day, now DST is inactive:
    fmt.Println(d1)
    // 2021-10-31 20:00:00 +0100 CET
}
```

---

### Credits
[colm.anseo](https://stackoverflow.com/users/1218512/colm-anseo) and [MrFuppes](https://stackoverflow.com/users/10197418/mrfuppes) for providing the following answers on Stackoverflow:  
* https://stackoverflow.com/a/68938947/175071
* https://stackoverflow.com/a/68966317/175071
