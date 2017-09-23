package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Schema	string `json:"schema"`
}

// GetPersistenceConfig generates the persistence struct from persistence.json
func GetPersistenceConfig() Persistence {
	raw, err := ioutil.ReadFile("./persistence.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var persistence Persistence

	json.Unmarshal(raw, &persistence)
	return persistence
}
