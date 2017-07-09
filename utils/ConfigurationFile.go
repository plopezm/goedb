package utils

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"goedb"
)

func GetPersistenceConfig() goedb.Persistence {
	raw, err := ioutil.ReadFile("./persistence.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var persistence goedb.Persistence

	json.Unmarshal(raw, &persistence)
	return persistence
}
