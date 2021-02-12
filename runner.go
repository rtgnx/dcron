package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Runner runs JobSpec on docker
type Runner client.Client

func (r Runner) createContainer(job JobSpec) (string, error) {
	c := client.Client(r)

	log.Printf("Creating Container for Job %s", job.Name)

	container, err := c.ContainerCreate(
		context.Background(), &container.Config{}, &container.HostConfig{}, nil, nil, job.dockerSafeName())

	return container.ID, err
}

func (r Runner) startContainer(id string) error {
	c := client.Client(r)
	return c.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func (r Runner) containerWait(id string) {
	c := client.Client(r)
	c.ContainerWait(context.Background(), id, container.WaitConditionNextExit)
}

func (r Runner) removeContainer(id string, force bool) error {
	c := client.Client(r)

	return c.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         force,
	})
}

// Run JobSpec
func (r Runner) Run(job JobSpec) error {

	id, err := r.createContainer(job)

	if err != nil {
		return err
	}

	defer func() {

		log.Printf("Removing container after job %s", job.Name)
		if err := r.removeContainer(id, false); err != nil {
			log.Printf("Job %s has finnished with errors", job.Name)
			log.Println(err.Error())
			return
		}
		log.Printf("Job %s has finished", job.Name)
	}()

	log.Printf("Starting job %s", job.Name)
	r.startContainer(id)
	r.containerWait(id)

	return nil
}
