package goedb

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/plopezm/goedb/database/dialect"

	"github.com/plopezm/goedb/config"
	"github.com/plopezm/goedb/database"
)

var goedbStandalone *dbm

type dbm struct {
	drivers map[string]database.EntityManager
}

func init() {
	log.Println("[GOEDB] library version: 1.0.0")
	goedbStandalone = new(dbm)
	goedbStandalone.drivers = make(map[string]database.EntityManager)
}

// Initialize gets the datasources from persistence.json
func Initialize() {
	var persistence config.Persistence
	persistence = config.GetPersistenceConfig("persistence.json")

	for _, datasource := range persistence.Datasources {
		driver := new(database.SQLDatabase)
		driver.Dialect = dialect.GetDialect(datasource.Driver)
		err := driver.Open(datasource.Driver, datasource.URL, datasource.Schema)
		if err != nil {
			fmt.Fprintf(os.Stdout, "[Connection ERROR for Persistence unit { %s } URL { %s }]: %v\n", datasource.Name, datasource.URL, err)
			//os.Exit(1)
			continue
		}
		goedbStandalone.drivers[datasource.Name] = driver
	}
}

// GetEntityManager returns a entity manager for the datasource selected.
func GetEntityManager(persistenceUnit string) (database.EntityManager, error) {
	entityManager, ok := goedbStandalone.drivers[persistenceUnit]
	if !ok {
		return nil, errors.New("Persistence unit \"" + persistenceUnit + "\" not found in persistence.json")
	}

	return entityManager, nil
}
