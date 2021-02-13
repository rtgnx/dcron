package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	job := new(JobSpec)
	err := job.fromFile("./manifests/hello_world.yaml")
	assert.Equal(t, nil, err)

	scheduler.Jobs = append(scheduler.Jobs, *job)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/jobs", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, getJobs(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		jobs := make([]JobSpec, 0)

		err := json.Unmarshal(rec.Body.Bytes(), &jobs)

		if err != nil {
			log.Fatalln(err.Error())
		}

		assert.Equal(t, err, nil)

		assert.Equal(t, job.Name, jobs[0].Name)
	}
}
