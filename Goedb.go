package goedb

import (
	"errors"
	"fmt"
	"github.com/plopezm/goedb/config"
	"github.com/plopezm/goedb/manager"
	"os"
	"github.com/plopezm/goedb/dialect"
)

var goedbStandalone *dbm

type dbm struct {
	drivers map[string]manager.EntityManager
}

func init() {
	var persistence config.Persistence
	goedbStandalone = new(dbm)
	persistence = config.GetPersistenceConfig()

	goedbStandalone.drivers = make(map[string]manager.EntityManager)
	for _, datasource := range persistence.Datasources {
		driver := new(manager.GoedbSQLDriver)
		err := driver.Open(datasource.Driver, datasource.URL, datasource.Schema)
		driver.Dialect = dialect.GetDialect(datasource.Driver)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
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
