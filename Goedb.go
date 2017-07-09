package goedb

import (
	"goedb/manager"
	"goedb/config"
	"errors"
)


var goedbStandalone *DBM


type DBM struct {
	drivers	map[string]manager.EntityManager
}



func init(){
	var persistence config.Persistence
	goedbStandalone = new(DBM);
	persistence = config.GetPersistenceConfig()

	goedbStandalone.drivers = make(map[string]manager.EntityManager)
	for _, datasource := range persistence.Datasources{
		driver := new(manager.GoedbSQLDriver)
		driver.Open(datasource.Driver, datasource.Url)
		goedbStandalone.drivers[datasource.Name] = driver
	}

}

func GetEntityManager(persistenceUnit string) (manager.EntityManager, error) {
	entityManager, ok := goedbStandalone.drivers[persistenceUnit]
	if !ok {
		return nil, errors.New("Persistence unit not found in persistence.json")
	}
	return entityManager, nil;
}





