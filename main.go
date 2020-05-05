package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"rickmaster2008/dm/handlers"
)

// RefreshRate used to set refresh rate
const RefreshRate = time.Millisecond * 100

func main() {
	url := os.Args[len(os.Args)-1]
	splits := strings.Split(url, "/")

	var n string
	flag.StringVar(&n, "name", splits[len(splits)-1], "Specify a different filename")
	flag.Parse()

	r, fsize := handlers.Head(url)

	if r {
		err := handlers.DownloadFile(n, url, fsize)
		if err != nil {
			panic(err)
		}
	} else {
		panic("El url no acepta Range")
	}
}
