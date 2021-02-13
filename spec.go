package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// JobSpec for cron jobs
type JobSpec struct {
	Name        string            `yaml:"name" json:"name"`
	Image       string            `yaml:"image" json:"image"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Commands    []string          `yaml:"commands" json:"commands,omitempty"`
	CronExpr    string            `yaml:"cron" json:"cron_expr"`
}

// Job type definition
type Job struct {
	Spec    JobSpec
	History []string
	Logs    map[string]io.Reader
}

func (job JobSpec) envPairs() []string {
	env := make([]string, 0)
	for k, v := range job.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

func (job *JobSpec) fromFile(filePath string) error {
	fd, err := os.Open(filePath)

	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(fd)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, job)
}
func (job *JobSpec) dockerSafeName() string {
	return strings.Replace(job.Name, " ", "_", -1)
}
