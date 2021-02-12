package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// JobSpec for cron jobs
type JobSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Commands    []string          `yaml:"commands"`
	CronExpr    string            `yaml:"cron"`
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
