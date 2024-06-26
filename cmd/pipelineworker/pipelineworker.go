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

  srv := asynq.NewServer(
    asynq.RedisClientOpt{
    	Addr: opts.Addr,
    	Username: opts.Username,
    	Password: opts.Password,
    	DB: opts.DB,
    },
    asynq.Config{
        Concurrency: 10,
    },
  )

  // mux maps a type to a handler
  mux := asynq.NewServeMux()
  mux.HandleFunc(tasks.TypeProcessImage, tasks.HandleProcessImageTask)

  if err := srv.Run(mux); err != nil {
      log.Fatalf("could not run server: %v", err)
  }
}
