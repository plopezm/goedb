package goedb

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"github.com/plopezm/goedb/manager"
	"github.com/stretchr/testify/assert"
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
	em, err = GetEntityManager("testSQLite3")
	if err != nil {
		panic("Persistence unit not defined in persistence.json")
	}
}

func TestD_ComplexStructs_Migrate(t *testing.T) {
	em, err := GetEntityManager("testSQLite3")
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&TestSoldier{})
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&TestTroop{})
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Complex_Model(t *testing.T) {
	soldier1, err := em.Model(&TestSoldier{})
	assert.Nil(t, err)
	assert.Equal(t, "TestSoldier", soldier1.Name)

	troop1, err := em.Model(&TestTroop{})
	assert.Nil(t, err)
	assert.Equal(t, "TestTroop", troop1.Name)
}

func TestDB_Complex_Insert(t *testing.T) {
	em, err := GetEntityManager("testSQLite3")
	if err != nil {
		t.Error(err)
	}

	troop1 := &TestTroop{
		Name: "TheBestTeam",
	}

	soldier1 := &TestSoldier{
		Name: "Ryan",
		Troop: *troop1,
	}

	_, err = em.Insert(troop1)
	assert.Nil(t, err)

	soldier1.Troop.Id = 1

	_, err = em.Insert(soldier1)
	assert.Nil(t, err)
}
