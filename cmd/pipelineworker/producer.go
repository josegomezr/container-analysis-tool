package main

import (
    "log"
    "os"

		"github.com/redis/go-redis/v9"
    "github.com/hibiken/asynq"
    "github.com/josegomezr/container-analysis-tools/pkg/tasks"
)

func main() {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "redis://127.0.0.1:6379/0"
	}

	opts, err := redis.ParseURL(redisUrl)

	if err != nil {
		log.Fatalf("Error parsing redis url: %v", err)
		return
	}

  client := asynq.NewClient(
    asynq.RedisClientOpt{
    	Addr: opts.Addr,
    	Username: opts.Username,
    	Password: opts.Password,
    	DB: opts.DB,
    },
  )
  defer client.Close()

  task, err := tasks.NewProcessImageTask("registry.suse.com/bci/bci-busybox")

  if err != nil {
      log.Fatalf("could not create task: %v", err)
  }
  info, err := client.Enqueue(task)
  if err != nil {
      log.Fatalf("could not enqueue task: %v", err)
  }
  log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
