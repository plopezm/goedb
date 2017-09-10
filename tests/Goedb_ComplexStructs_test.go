package tests

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"github.com/plopezm/goedb/manager"
	"github.com/stretchr/testify/assert"
	"github.com/plopezm/goedb"
	"database/sql"
)

type TestTroop struct {
	Id		int		`goedb:"pk,autoincrement"`
	Name	string		`goedb:"unique"`
}

type TestSoldier struct {
	Id    	int 		`goedb:"pk,autoincrement"`
	Name    string		`goedb:"unique"`
	Troop	TestTroop	`goedb:"fk=TestTroop(Id)"`
}

var em manager.EntityManager

func init(){
	var err error
	em, err = goedb.GetEntityManager("testSQLite3")
	if err != nil {
		panic("Persistence unit not defined in persistence.json")
	}
}

func Test_Goedb_Migrate(t *testing.T) {
	err := em.Migrate(&TestSoldier{})
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&TestTroop{})
	if err != nil {
		t.Error(err)
	}
}

func Test_Goedb_Model(t *testing.T) {
	soldier1, err := em.Model(&TestSoldier{})
	assert.Nil(t, err)
	assert.Equal(t, "TestSoldier", soldier1.Name)

	troop1, err := em.Model(&TestTroop{})
	assert.Nil(t, err)
	assert.Equal(t, "TestTroop", troop1.Name)
}

func Test_Goedb_Insert(t *testing.T) {
	troop1 := &TestTroop{
		Name: "TheBestTeam",
	}

	soldier1 := &TestSoldier{
		Name: "Ryan",
		Troop: *troop1,
	}

	_, err := em.Insert(troop1)
	assert.Nil(t, err)

	soldier1.Troop.Id = 1

	_, err = em.Insert(soldier1)
	assert.Nil(t, err)

	soldier1.Troop.Id = 2
	_, err = em.Insert(soldier1)
	assert.NotNil(t, err)
}

func Test_Goedb_First_By_Id(t *testing.T){
	soldier1 := &TestSoldier{
		Id: 1,
	}

	err := em.First(soldier1, "")
	assert.Nil(t, err)
	assert.Equal(t, "Ryan", soldier1.Name)
	assert.Equal(t, 1, soldier1.Troop.Id)
}

func Test_Goedb_First_By_Name(t *testing.T){
	soldier1 := &TestSoldier{
		Name: "Ryan",
	}

	err := em.First(soldier1, "TestSoldier.Name = :name", sql.Named("name", "Ryan"))
	assert.Nil(t, err)
	assert.Equal(t, 1, soldier1.Id)
	assert.Equal(t, 1, soldier1.Troop.Id)
}
