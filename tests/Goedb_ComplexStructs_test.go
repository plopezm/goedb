package tests

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/plopezm/goedb"
	"github.com/plopezm/goedb/manager"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTroop struct {
	Id   int    `goedb:"pk,autoincrement"`
	Name string `goedb:"unique"`
}

type TestSoldier struct {
	Id    int       `goedb:"pk,autoincrement"`
	Name  string    `goedb:"unique"`
	Troop TestTroop `goedb:"fk=TestTroop(Id)"`
}

var em manager.EntityManager

func init() {
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
		Name:  "Ryan",
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

func Test_Goedb_First_By_PrimaryKey(t *testing.T) {
	soldier1 := &TestSoldier{
		Id: 1,
	}

	err := em.First(soldier1, "")
	assert.Nil(t, err)
	assert.Equal(t, "Ryan", soldier1.Name)
	assert.Equal(t, 1, soldier1.Troop.Id)
}

func Test_Goedb_First_By_Name(t *testing.T) {
	soldier1 := &TestSoldier{
		Name: "Ryan",
	}

	err := em.First(soldier1, "TestSoldier.Name = :name", sql.Named("name", "Ryan"))
	assert.Nil(t, err)
	assert.Equal(t, 1, soldier1.Id)
	assert.Equal(t, 1, soldier1.Troop.Id)
}

func weapon_call() (*TestSoldier, *TestSoldier, *TestSoldier, *TestSoldier) {
	soldier2 := &TestSoldier{
		Name: "Bryan",
		Troop: TestTroop{
			Id: 1,
		},
	}
	soldier3 := &TestSoldier{
		Name: "Steve",
		Troop: TestTroop{
			Id: 1,
		},
	}
	soldier4 := &TestSoldier{
		Name: "Eduard",
		Troop: TestTroop{
			Id: 1,
		},
	}
	soldier5 := &TestSoldier{
		Name: "Chuck",
		Troop: TestTroop{
			Id: 1,
		},
	}
	return soldier2, soldier3, soldier4, soldier5
}

func Test_Find_All_Soldiers(t *testing.T) {
	s1, s2, s3, s4 := weapon_call()
	em.Insert(s1)
	em.Insert(s2)
	em.Insert(s3)
	em.Insert(s4)

	foundSoldiers := make([]TestSoldier, 0)

	err := em.Find(&foundSoldiers, "TestTroop.Id = :troop_id", sql.Named("troop_id", 1))
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 5, len(foundSoldiers))
}

func Test_Find_One_Soldier(t *testing.T) {
	foundSoldiers := make([]TestSoldier, 0)
	err := em.Find(&foundSoldiers, "TestSoldier.Id = :soldier_id", sql.Named("soldier_id", 3))
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 1, len(foundSoldiers))
}

func Test_Delete_Soldier_By_PrimaryKey(t *testing.T) {
	soldier1 := &TestSoldier{
		Id: 3,
	}
	result, err := em.Remove(soldier1, "")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.NumRecordsAffected)
}

func Test_Delete_Soldier_By_OtherField(t *testing.T) {
	soldier1 := &TestSoldier{}
	result, err := em.Remove(soldier1, "TestSoldier.Name = :name", sql.Named("name", "Bryan"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.NumRecordsAffected)
}

func Test_DropTable(t *testing.T) {
	err := em.DropTable(&TestTroop{})
	assert.Nil(t, err)
	err = em.DropTable(&TestSoldier{})
	assert.Nil(t, err)
	_, err = em.Model(&TestTroop{})
	assert.NotNil(t, err)
	_, err = em.Model(&TestSoldier{})
	assert.NotNil(t, err)
}
