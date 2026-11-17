## Get the full name of your local time zone (OS setting)

Works on windows, linux and macos

### Why?
built-in functionality sometimes won't suffice:
```go
zone, _ := time.Now().Zone() // my time zone
fmt.Println(zone) // e.g. "CEST"

loc, err := time.LoadLocation(zone)
if err != nil {
    // the program panics here:
    panic("fail: could not load location from zone")
}

// unreachable:
fmt.Println("success: location loaded from zone")
```
localizing a date with obtained `loc` will cause an error because it's not in the IANA format.
> panic: time: missing Location in call to Date

While tzlocal gives the correct IANA name
```
tzname, _ := tzlocal.RuntimeTZ() // assuming error is handled
fmt.Printf("Actual IANA name: %v\n", tzname) // Prints "Actual IANA name: Europe/Paris"
```

---

### Package Usage
```
go get github.com/thlib/go-timezone-local/tzlocal
```

### Included example

A runnable example is available in [`example/main.go`](./example/main.go). It prints the operating system, detected IANA time zone, and current local time.

The example imports [`time/tzdata`](https://pkg.go.dev/time/tzdata), so the IANA time zone database is embedded in the executable. This is useful for standalone Windows builds where a time zone database may not otherwise be available.

Run it from the repository root:

```sh
go run ./example
```

Build and run it on Linux:

```sh
go build -o bin/timezone-example-linux ./example
./bin/timezone-example-linux
```

Cross-compile it for 64-bit Windows:

```sh
GOOS=windows GOARCH=amd64 go build -o bin/timezone-example-windows-amd64.exe ./example
```

Then run it from PowerShell on Windows:

```powershell
.\bin\timezone-example-windows-amd64.exe
```

The generated `bin/` directory is ignored by Git.

### Use in your own program

Open your project folder  
Create a file `main.go`

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

Run the following commands:
```
go mod init example.com/yourpackage
go mod vendor
go run main.go
```

It should print the go runtime timezone.


### For contributors to update the Windows time zone mapping

The generator downloads the latest IANA legacy names and Unicode CLDR Windows
zone mappings, then rewrites `tzlocal/tzmapping.go`. It therefore requires
network access to `data.iana.org` and `raw.githubusercontent.com`.

From the repository root, run:

```sh
go generate
```

Generation compares the actual key/value pairs in each table and leaves the file
untouched when they match; dates and formatting do not count as changes. An update
is blocked if either table loses any existing key or has less than 90% similarity.
Similarity is the number of identical key/value pairs divided by the larger of
the old and new table sizes, calculated separately for each table. Both tables
must be nonempty, all IANA names must load, and references between tables must exist.
A blocked update leaves the existing file intact and prints the reason.

The [Time zone mappings workflow](.github/workflows/tzmapping.yml) checks upstream
daily at 07:23 UTC and supports manual runs from the Actions tab on `main`. It
tests the candidate on Linux, macOS, and Windows before committing the mapping
directly to `main`, atomically pushing a new minor version tag, and publishing a
GitHub release. For example, `v0.0.9` becomes `v0.1.0`; unchanged mappings produce
no commit or new version. Failed runs open an issue assigned to `timohuovinen`,
or update the existing open notification issue with the latest logs. A rerun
retries a release if the validated commit and tag were already pushed.

The workflow uses the built-in `GITHUB_TOKEN` with contents/issue write permissions;
no extra secret is needed. Repository rules must permit its direct commits to
`main` and version tags. If a rule blocks publication, the run fails and notifies
the maintainer; it does not bypass the rule. GitHub may disable scheduled workflows
in public repositories after 60 days without repository activity; re-enable the
workflow in Actions if that happens. Delivery timing follows GitHub's scheduler.

To approve an intentional removal or a larger upstream change, review and commit
the mapping manually, run the tests, and create its release manually. The scheduled
workflow has no switch to bypass its safeguards.

If it reports `GOPROXY list is not the empty string, but contains no entries`, the
configured `GOPROXY` is malformed. 
To override it for this command in a POSIX shell or WSL, run:

```sh
GOPROXY=https://proxy.golang.org,direct go generate
```

### Credits

All credit goes to the various contributors, [colm.anseo](https://stackoverflow.com/users/1218512/colm-anseo) and [MrFuppes](https://stackoverflow.com/users/10197418/mrfuppes) for providing the following answers:  
* https://stackoverflow.com/a/68938947/175071
* https://stackoverflow.com/a/68966317/175071

## License

This project is **dual‑licensed**:  
- **The Unlicense** (public domain). See [LICENSE](./LICENSE).  
- **Apache License 2.0**. See [LICENSE-APACHE.txt](./LICENSE-APACHE.txt).  

You may choose **either** license at your option.

### Third‑party licences
>This repository depends on several third‑party Go packages, for example `golang.org/x/sys` that are distributed under the BSD 3‑Clause licence.
