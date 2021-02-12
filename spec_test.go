package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvPairs(t *testing.T) {

	spec := new(JobSpec)
	spec.Environment = make(map[string]string)
	spec.Environment["KEY"] = "value"

	pairs := spec.envPairs()

	assert.Equal(t, len(pairs), 1, "Should have one element")
	assert.Equal(t, pairs[0], "KEY=value")
}
