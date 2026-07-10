package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
	_ "time/tzdata"

	"github.com/thlib/go-timezone-local/tzlocal"
)

func main() {
	timeZone, err := tzlocal.RuntimeTZ()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load time zone %q: %v\n", timeZone, err)
		os.Exit(1)
	}

	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Time zone: %s\n", timeZone)
	fmt.Printf("Local time: %s\n", time.Now().In(location).Format(time.RFC3339))
}
