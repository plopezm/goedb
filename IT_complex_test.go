package goedb

import (
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type TestTroop struct {
	ID       int           `goedb:"pk,autoincrement"`
	Name     string        `goedb:"unique"`
	Soldiers []TestSoldier `goedb:"ignore"`
}

type TestSoldier struct {
	ID    int       `goedb:"pk,autoincrement"`
	Name  string    `goedb:"unique"`
	Troop TestTroop `goedb:"fk=TestTroop(ID)"`
}

type TestCustomSoldier struct {
	ID        int
	Name      string
	TroopName string
}

const persistenceUnitItComplexTest = "testSQLite3"

func init() {
	Initialize()
}

func Test_Goedb_Open_Fail(t *testing.T) {
	_, err := GetEntityManager("Not-exists")
	assert.NotNil(t, err)
}

func Test_Goedb_Open_And_Close(t *testing.T) {
	em, err := GetEntityManager("closeTest")
	assert.Nil(t, err)

	tx, err := em.TxBegin()
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	tx.Rollback()

	err = em.Close()
	assert.Nil(t, err)
}

func Test_Goedb_Migrate(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.Migrate(&TestTroop{}, true, true)
	assert.Nil(t, err)

	err = em.Migrate(&TestSoldier{}, true, true)
	assert.Nil(t, err)
}

func Test_Goedb_Migrate_Recreate(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.Migrate(&TestTroop{}, true, false)
	assert.NotNil(t, err)

	err = em.Migrate(&TestSoldier{}, true, false)
	assert.NotNil(t, err)
}

func Test_Goedb_Model(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1, err := em.Model(&TestSoldier{})
	assert.Nil(t, err)
	assert.Equal(t, "TestSoldier", soldier1.Name)

	troop1, err := em.Model(&TestTroop{})
	assert.Nil(t, err)
	assert.Equal(t, "TestTroop", troop1.Name)
}

func Test_Goedb_Insert(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	troop1 := &TestTroop{
		Name: "TheBestTeam",
	}

	soldier1 := &TestSoldier{
		Name:  "Ryan",
		Troop: *troop1,
	}

	_, err = em.Insert(troop1)
	assert.Nil(t, err)

	soldier1.Troop.ID = 1

	_, err = em.Insert(soldier1)
	assert.Nil(t, err)

	soldier1.Troop.ID = 2
	_, err = em.Insert(soldier1)
	assert.NotNil(t, err)
}

func Test_Goedb_First_By_PrimaryKey(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1 := &TestSoldier{
		ID: 1,
	}

	err = em.First(soldier1, "", nil)
	assert.Nil(t, err)
	assert.Equal(t, "Ryan", soldier1.Name)
	assert.Equal(t, 1, soldier1.Troop.ID)
}

func Test_Goedb_First_By_Name(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1 := &TestSoldier{
		Name: "Ryan",
	}

	err = em.First(soldier1, "TestSoldier.Name = :name", map[string]interface{}{"name": "Ryan"})
	assert.Nil(t, err)
	assert.Equal(t, 1, soldier1.ID)
	assert.Equal(t, 1, soldier1.Troop.ID)
}

func Test_Goedb_Native_First_By_Name(t *testing.T) {

	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	var customSoldier TestCustomSoldier

	err = em.NativeFirst(&customSoldier, "SELECT ts.ID, ts.Name, tt.Name FROM TestSoldier ts, TestTroop tt WHERE ts.Name = :name AND ts.Troop = tt.ID", map[string]interface{}{"name": "Ryan"})
	assert.Nil(t, err)
	assert.Equal(t, 1, customSoldier.ID)
	assert.Equal(t, "TheBestTeam", customSoldier.TroopName)
}

func weaponCall() (*TestSoldier, *TestSoldier, *TestSoldier, *TestSoldier) {
	soldier2 := &TestSoldier{
		Name: "Bryan",
		Troop: TestTroop{
			ID:   1,
			Name: "ExampleTroop",
		},
	}
	soldier3 := &TestSoldier{
		Name: "Steve",
		Troop: TestTroop{
			ID:   1,
			Name: "ExampleTroop",
		},
	}
	soldier4 := &TestSoldier{
		Name: "Eduard",
		Troop: TestTroop{
			ID:   1,
			Name: "ExampleTroop",
		},
	}
	soldier5 := &TestSoldier{
		Name: "Chuck",
		Troop: TestTroop{
			ID:   1,
			Name: "ExampleTroop",
		},
	}
	return soldier2, soldier3, soldier4, soldier5
}

func Test_Find_All_Soldiers(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	s1, s2, s3, s4 := weaponCall()
	em.Insert(s1)
	em.Insert(s2)
	em.Insert(s3)
	em.Insert(s4)

	foundSoldiers := make([]TestSoldier, 0)

	err = em.Find(&foundSoldiers, "", nil)
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 5, len(foundSoldiers))
}

func Test_Native_Find_All_Soldiers(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	s1, s2, s3, s4 := weaponCall()
	em.Insert(s1)
	em.Insert(s2)
	em.Insert(s3)
	em.Insert(s4)

	foundSoldiers := make([]TestCustomSoldier, 0)

	err = em.NativeFind(&foundSoldiers, "SELECT ts.ID, ts.Name, tt.Name FROM TestSoldier ts, TestTroop tt WHERE ts.Troop = tt.ID", nil)
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 5, len(foundSoldiers))
}

func Test_Find_One_Soldier(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	foundSoldiers := make([]TestSoldier, 0)
	err = em.Find(&foundSoldiers, "TestSoldier.Name = :soldier_name", map[string]interface{}{"soldier_name": "Steve"})
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 1, len(foundSoldiers))
	assert.Equal(t, "Steve", foundSoldiers[0].Name)
}

func Test_Native_Find_One_Soldier(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	foundSoldiers := make([]TestCustomSoldier, 0)
	err = em.NativeFind(&foundSoldiers, "SELECT ts.ID, ts.Name, tt.Name FROM TestSoldier ts, TestTroop tt WHERE ts.Name = :soldier_name AND ts.Troop = tt.ID", map[string]interface{}{"soldier_name": "Steve"})
	assert.Nil(t, err)
	assert.NotNil(t, foundSoldiers)
	assert.Equal(t, 1, len(foundSoldiers))
	assert.Equal(t, "Steve", foundSoldiers[0].Name)
}

func Test_Update_Soldier_By_PrimaryKey(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1 := &TestSoldier{
		ID:   3,
		Name: "UpdateTest",
		Troop: TestTroop{
			ID: 1,
		},
	}

	result, err := em.Update(soldier1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.NumRecordsAffected)

	soldier1.Name = ""
	err = em.First(soldier1, "", nil)
	assert.Nil(t, err)
	assert.Equal(t, "UpdateTest", soldier1.Name)
	assert.Equal(t, 1, soldier1.Troop.ID)
}

func Test_Delete_Soldier_By_PrimaryKey(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1 := &TestSoldier{
		ID: 3,
	}
	result, err := em.Remove(soldier1, "", nil)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.NumRecordsAffected)
}

func Test_Delete_Soldier_By_OtherField(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	soldier1 := &TestSoldier{}
	result, err := em.Remove(soldier1, "TestSoldier.Name = :soldier_name", map[string]interface{}{"soldier_name": "Chuck"})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result.NumRecordsAffected)
}

func Test_DropTable(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItComplexTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.DropTable(&TestSoldier{})
	assert.Nil(t, err)
	err = em.DropTable(&TestTroop{})
	assert.Nil(t, err)
	_, err = em.Model(&TestTroop{})
	assert.NotNil(t, err)
	_, err = em.Model(&TestSoldier{})
	assert.NotNil(t, err)
}
