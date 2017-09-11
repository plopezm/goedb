package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Persistence struct {
	Datasources []Datasource
}

type Datasource struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Url    string `json:"url"`
}

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
