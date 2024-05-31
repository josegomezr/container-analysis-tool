package main

import (
	"context"
	"fmt"
	"github.com/josegomezr/container-analysis-tools/pkg/tasks"
	"os"
)

func main() {
	ctx := context.Background()
	src := os.Args[1]
	dest := os.Args[2]

	if err := os.MkdirAll(dest, 0755); err != nil {
		fmt.Printf("err: %v \n", err)
		return
	}

	tasks.ProcessImage(ctx, src, dest)
}
