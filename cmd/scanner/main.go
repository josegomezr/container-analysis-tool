package main

import (
	"context"
	"fmt"
	"github.com/josegomezr/container-analysis-tools/pkg/scanner"
	"os"
)

func main() {
	ctx := context.TODO()
	fmt.Printf("START \n")

	image := os.Args[1]
	scan, err := scanner.ScanImage(ctx, image)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Println(scan)
}
