package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPersistenceConfig(t *testing.T) {
	persistence := GetPersistenceConfig("persistence.json")
	assert.Equal(t, 2, len(persistence.Datasources))
}

func TestGetPersistenceConfigNotFound(t *testing.T) {
	persistence := GetPersistenceConfig("Notfoundfile.json")
	assert.Equal(t, 0, len(persistence.Datasources))
}
