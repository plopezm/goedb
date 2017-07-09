package goedb

import (
	"goedb/drivers"
	"goedb/utils"
)


var goedbStandalone *DBM


type DBM struct {
	drivers	map[string]drivers.GoedbDriver

}

type Persistence struct{
	datasources		 []Datasource
}

type Datasource struct{
	Name	string  	`json:"name"`
	Driver 	string		`json:"driver"`
	Url		string		`json:"url"`
}



func init(){
	goedbStandalone = new(DBM);
	persistence := utils.GetPersistenceConfig()

	goedbStandalone.drivers = make(map[string]drivers.GoedbDriver)
	for _, datasource := range persistence.datasources{
		var driver = &drivers.GoedbSQLDriver{}
		driver.Open(datasource.Driver, datasource.Url)
		goedbStandalone.drivers[datasource.Name] = driver
	}

}


func GetInstance() *DBM {
	return goedbStandalone;
}

/*
func (dbm *DBM) SetDriver(driver drivers.GoedbDriver){
	dbm.driver = driver
}
*/

/*
func (dbm *DBM) Open(driver string, params string) error{
	return dbm.driver.Open(driver, params)
}
*/

func (dbm *DBM) Close() error{
	return dbm.driver.Close()
}

func (dbm *DBM) Migrate(i interface{}) (error){
	return dbm.driver.Migrate(i)
}

func (dbm *DBM) Model(i interface{}) (drivers.GoedbTable, error){
	return dbm.driver.Model(i)
}

func (dbm *DBM) Insert(i interface{})(drivers.GoedbResult, error){
	return dbm.driver.Insert(i)
}

func (dbm *DBM) Remove(i interface{})(drivers.GoedbResult, error){
	return dbm.driver.Remove(i)
}

func (dbm *DBM) First(i interface{}, where string) (error){
	return dbm.driver.First(i, where)
}

func (dbm *DBM) Find(i interface{}, where string) error{
	return dbm.driver.Find(i, where)
}

func (dbm *DBM) DropTable(i interface{}) error{
	return dbm.driver.DropTable(i)
}



