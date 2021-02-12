package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
)

var (
	app       = cli.App("dcron", "Dockerized Cron")
	scheduler Scheduler
	e         = echo.New()
)

func main() {
	app.Command("run", "run dcron", cmdRun)
	if err := scheduler.Init(); err != nil {
		log.Fatalln(err.Error())
	}
	app.Run(os.Args)

	e.GET("/api/jobs", getJobs)
	e.Start(":80")
}

func cmdRun(cmd *cli.Cmd) {
	var (
		manifests = cmd.StringOpt("manifests", "./manifests", "Path to manifests directory")
	)

	// Run this function when the command is invoked
	cmd.Action = func() {
		cwd, _ := os.Getwd()
		_, err := os.Stat(path.Join(cwd, *manifests))

		if os.IsNotExist(err) {
		}

		fileList, _ := ioutil.ReadDir(path.Join(cwd, *manifests))

		for _, info := range fileList {
			spec := new(JobSpec)
			if err := spec.fromFile(path.Join(cwd, *manifests, info.Name())); err != nil {
				log.Fatalln(err.Error())
			}
			scheduler.Jobs = append(scheduler.Jobs, *spec)
		}

		scheduler.Start()

	}
}
