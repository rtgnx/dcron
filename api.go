package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func postJobSpec(ctx echo.Context) error {
	job := new(JobSpec)

	if err := ctx.Bind(job); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := scheduler.scheduleJob(*job); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	scheduler.Jobs = append(scheduler.Jobs, *job)

	return ctx.NoContent(http.StatusCreated)
}

func getJobs(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, scheduler.Jobs)
}
