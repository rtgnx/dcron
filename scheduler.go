package main

import (
	"io"
	"log"
	"sync"

	"github.com/docker/docker/client"
	"github.com/robfig/cron/v3"
)

// Scheduler for cron jobs
type Scheduler struct {
	Jobs   []JobSpec
	Logs   map[string]io.Reader
	cron   *cron.Cron
	wg     sync.WaitGroup
	runner Runner
}

// Init Scheduler
func (scheduler *Scheduler) Init() error {
	scheduler.Logs = make(map[string]io.Reader)
	scheduler.cron = cron.New(cron.WithSeconds())
	cli, err := client.NewEnvClient()

	if err != nil {
		return err
	}

	scheduler.runner = Runner(*cli)
	return nil
}

// Start Scheduler
func (scheduler *Scheduler) Start() {

	scheduler.cron.Start()

	for _, job := range scheduler.Jobs {
		log.Printf("Scheduling %s [%s]", job.Name, job.CronExpr)
		_, err := scheduler.cron.AddFunc(job.CronExpr, func() {

			if err := scheduler.runner.Run(job); err != nil {
				log.Println(err.Error())
			}

			scheduler.wg.Done()
		})

		if err != nil {
			log.Println(err.Error())
			continue
		}

		scheduler.wg.Add(1)
	}

	defer scheduler.cron.Stop()
	scheduler.wg.Wait()
}
