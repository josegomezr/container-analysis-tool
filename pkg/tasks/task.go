package tasks

import (
	"context"
	"fmt"
	"github.com/josegomezr/container-analysis-tools/pkg/downloader"
	"github.com/josegomezr/container-analysis-tools/pkg/scanner"
	"os"
	"sync"
)

func doScan(ctx context.Context, inputFile string) {
	scan, err := scanner.ScanImage(ctx, inputFile)

	if err != nil {
		fmt.Printf("err scanning: %+v\n", err)
		return
	}

	if len(scan) > 0 {
		fmt.Printf("") // I hate this variable defined-not-used...
	}

	// TODO: replace with the real stuff
	fmt.Printf("curl https://my-endpoint/scanner-result -d data=%v\n", len(scan))
}

func doStat(ctx context.Context, inputFile string) {
	finfo, err := os.Stat(inputFile)

	if err != nil {
		fmt.Printf("err stat'ing: %+v\n", err)
		return
	}

	if finfo.Size() > 0 {
		fmt.Printf("") // I hate this variable defined-not-used...
	}

	// TODO: replace with the real stuff
	fmt.Printf("curl https://my-endpoint/size -d data=%v\n", finfo.Size())
}

func ProcessImage(ctx context.Context, src string, workingDir string) {
	fmt.Println("Start of task")

	downloadResults, err := downloader.Download(ctx, &downloader.Opts{
		Source:            src,
		DestinationFolder: workingDir,
	})

	if err != nil {
		fmt.Printf("error downloading: %v \n", err)
		return
	}

	var wg sync.WaitGroup

	for _, result := range downloadResults {
		wg.Add(1)
		go func(result downloader.DownloadResult) {
			defer wg.Done()
			doScan(ctx, result.Filename)
		}(result)

		wg.Add(1)
		go func(result downloader.DownloadResult) {
			defer wg.Done()
			doScan(ctx, result.Filename)
		}(result)
	}

	// wait for all stats & scans to complete
	wg.Wait()

	fmt.Println("Task Done")
}
