package tests

import (
	"testing"
	"github.com/plopezm/goedb/config"
	"github.com/stretchr/testify/assert"
)

func TestGetPersistenceConfig(t *testing.T){
	persistence := config.GetPersistenceConfig("persistence.json")
	assert.Equal(t,2, len(persistence.Datasources))
}

func TestGetPersistenceConfigNotFound(t *testing.T){
	persistence := config.GetPersistenceConfig("Notfoundfile.json")
	assert.Equal(t,0, len(persistence.Datasources))
}
