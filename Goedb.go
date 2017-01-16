package goedb

type DBM struct {
	driver GoedbDriver
}

func init(){

}

func (dbm *DBM) SetDriver(driver GoedbDriver){
	dbm.driver = driver
}

func (dbm *DBM) Open(driver string, params string) error{
	return dbm.driver.Open(driver, params)
}

func (dbm *DBM) Close() error{
	return dbm.driver.Close()
}

func (dbm *DBM) Migrate(i interface{}) (error){
	return dbm.driver.Migrate(i)
}

func (dbm *DBM) Model(i interface{}) (GoedbTable, error){
	return dbm.driver.Model(i)
}

func (dbm *DBM) Insert(i interface{})(GoedbResult, error){
	return dbm.driver.Insert(i)
}

func (dbm *DBM) Remove(i interface{})(GoedbResult, error){
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



