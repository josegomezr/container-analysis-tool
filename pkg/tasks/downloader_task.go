package tasks

import (
	"context"
	"fmt"
	"log"
	"github.com/josegomezr/container-analysis-tools/pkg/downloader"
	"os"
	"sync"
	"encoding/json"

  "github.com/hibiken/asynq"
)

const (
  TypeProcessImage   = "process-image"
)

type ProcessImagePayload struct {
    ImageName string
}

func NewProcessImageTask(imageName string) (*asynq.Task, error) {
    payload, err := json.Marshal(ProcessImagePayload{ImageName: imageName})

    if err != nil {
        return nil, err
    }

    return asynq.NewTask(TypeProcessImage, payload), nil
}

func HandleProcessImageTask(ctx context.Context, t *asynq.Task) error {
	log.Println("Start of task")

	var p ProcessImagePayload
  if err := json.Unmarshal(t.Payload(), &p); err != nil {
      return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
  }

  src := fmt.Sprintf("docker://%s", p.ImageName)
  workingDir, err := os.MkdirTemp("./image/", "image-download-*")
  if err != nil {
		return fmt.Errorf("os.MkdirTemp failed: %v: %w", err, asynq.SkipRetry)
	}

	downloadResults, err := downloader.Download(ctx, &downloader.Opts{
		Source:            src,
		DestinationFolder: workingDir,
	})

	if err != nil {
		return fmt.Errorf("downloader.Download failed: %v: %w", err, asynq.SkipRetry)
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
			doStat(ctx, result.Filename)
		}(result)
	}

	// wait for all stats & scans to complete
	wg.Wait()

	fmt.Println("Task Done")
	return nil
}
