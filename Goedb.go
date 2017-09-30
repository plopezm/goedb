package goedb

import (
	"errors"
	"fmt"
	"github.com/plopezm/goedb/config"
	"github.com/plopezm/goedb/dialect"
	"github.com/plopezm/goedb/manager"
	"os"
)

var goedbStandalone *dbm

type dbm struct {
	drivers map[string]manager.EntityManager
}

func init() {
	goedbStandalone = new(dbm)
	goedbStandalone.drivers = make(map[string]manager.EntityManager)
}

// Initialize gets the datasources from persistence.json
func Initialize() {
	var persistence config.Persistence
	persistence = config.GetPersistenceConfig()

	for _, datasource := range persistence.Datasources {
		driver := new(manager.GoedbSQLDriver)
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
func GetEntityManager(persistenceUnit string) (manager.EntityManager, error) {
	entityManager, ok := goedbStandalone.drivers[persistenceUnit]
	if !ok {
		return nil, errors.New("Persistence unit not found in persistence.json")
	}

	return entityManager, nil
}
