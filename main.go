package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
)

var (
	app       = cli.App("dcron", "Dockerized Cron")
	scheduler Scheduler
	e         = echo.New()
	s3Bucket  = os.Getenv("S3_BUCKET")
	s3Prefix  = os.Getenv("S3_PREFIX")
)

func main() {
	app.Command("run", "run dcron", cmdRun)
	if err := scheduler.Init(); err != nil {
		log.Fatalln(err.Error())
	}
	app.Run(os.Args)

	e.GET("/api/jobs", getJobs)
	e.Start(":8080")
}

func cmdRun(cmd *cli.Cmd) {
	var (
		manifests = cmd.StringOpt("manifests", "./manifests", "Path to manifests directory")
		useS3     = cmd.BoolOpt("use-s3", false, "Use S3 bucket as source of manifests file")
	)

	// Run this function when the command is invoked
	cmd.Action = func() {
		cwd, _ := os.Getwd()
		if *useS3 {
			fromObjectStorage()
		} else {
			fromLocalStorage(path.Join(cwd, *manifests))
		}
		scheduler.Start()
	}

}

func fromObjectStorage() (err error) {
	scheduler.Jobs, err = ReadJobSpecs(s3Bucket, s3Prefix)
	return
}

func fromLocalStorage(rootPath string) {
	_, err := os.Stat(rootPath)

	if os.IsNotExist(err) {
	}

	fileList, _ := ioutil.ReadDir(rootPath)

	for _, info := range fileList {
		spec := new(JobSpec)
		if err := spec.fromFile(path.Join(rootPath, info.Name())); err != nil {
			log.Fatalln(err.Error())
		}
		scheduler.Jobs = append(scheduler.Jobs, *spec)
	}
}
