package config

import (
	"encoding/json"
	"io/ioutil"
)

// Persistence represents the collection of datasources defined by the developer in persistence.json
type Persistence struct {
	Datasources []Datasource
}

// Datasource represents the metadata of a connection pool
type Datasource struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	URL    string `json:"url"`
	Schema string `json:"schema"`
}

// GetPersistenceConfig generates the persistence struct from persistence.json
func GetPersistenceConfig(persistenceConfigFile string) Persistence {
	var persistence Persistence
	raw, err := ioutil.ReadFile(persistenceConfigFile)
	if err == nil {
		json.Unmarshal(raw, &persistence)
	} else {
		persistence.Datasources = make([]Datasource, 0)
	}
	return persistence
}
