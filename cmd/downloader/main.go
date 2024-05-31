package main

import (
	"context"
	"fmt"
	"github.com/josegomezr/container-analysis-tools/pkg/downloader"
	"os"
)

func main() {
	ctx := context.TODO()
	fmt.Println("Downloading")

	src := os.Args[1]
	dest := os.Args[2]

	if err := os.MkdirAll(dest, 0755); err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	downloader.Download(ctx, &downloader.Opts{
		Source:            src,
		DestinationFolder: dest,
	})

	fmt.Println("Done")
}
