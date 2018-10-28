package goedb

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type testuser struct {
	Email    string `goedb:"pk"`
	Password string
	Role     string
	DNI      int `goedb:"unique"`
	Admin    bool
}

type testcompany struct {
	UserEmail string `goedb:"fk=testuser(Email)"`
	Name      string
	Cif       string `goedb:"pk"`
}

type testusercompany struct {
	Email string `goedb:"pk,fk=testuser(Email)"`
	Cif   string `goedb:"pk,fk=testcompany(Cif)"`
}

type otherstruct struct {
	Asd   string
	Other string
}

const persistenceUnitItTest = "testSQLite3"

func init() {
	Initialize()
}

func TestOpen(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	if err != nil {
		t.Error(err)
	}

	db := em.GetDBConnection()
	assert.NotNil(t, db)
}

func TestDB_Migrate(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.Migrate(&testuser{}, true, true)
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&testcompany{}, true, true)
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&testusercompany{}, true, true)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Model(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	user, _ := em.Model(&testuser{})
	if user.Name != "testuser" || len(user.Columns) == 0 {
		t.Error("Error getting db model")
	}

	company, _ := em.Model(&testcompany{})
	if company.Name != "testcompany" || len(company.Columns) == 0 {
		t.Error("Error getting db model")
	}
}

func TestDB_Model_Not_Found(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	_, err = em.Model(&otherstruct{})
	if err == nil {
		t.Error("The result must has a error because the struct was not created")
	}
}

func TestDB_Insert(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser1 := &testuser{
		Email:    "Plm",
		Password: "asd",
		Role:     "asd",
		DNI:      123,
		Admin:    true,
	}

	newUser2 := &testuser{
		Email:    "Plm2",
		Password: "asd",
		Role:     "asd",
		DNI:      1234,
		Admin:    true,
	}

	newUser3 := &testuser{
		Email:    "Plm3",
		Password: "asd",
		Role:     "asd",
		DNI:      1235,
		Admin:    false,
	}

	_, err = em.Insert(newUser1)
	if err != nil {
		t.Error(err)
	}

	_, err = em.Insert(newUser2)
	if err != nil {
		t.Error(err)
	}

	_, err = em.Insert(newUser3)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Insert_with_FKs(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newComp1 := &testcompany{
		UserEmail: "Plm",
		Name:      "asd",
		Cif:       "asd1",
	}

	newComp2 := &testcompany{
		UserEmail: "Plm",
		Name:      "asd",
		Cif:       "asd2",
	}

	_, err = em.Insert(newComp1)
	if err != nil {
		t.Error(err)
	}

	_, err = em.Insert(newComp2)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Insert_Constraints(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newComp1 := &testcompany{
		UserEmail: "Plm",
		Name:      "asd",
		Cif:       "asd1",
	}

	newComp2 := &testcompany{
		UserEmail: "Plm",
		Name:      "asd123",
		Cif:       "asd2",
	}

	newComp3 := &testcompany{
		UserEmail: "Plm23455",
		Name:      "asd",
		Cif:       "asd4",
	}

	_, err = em.Insert(newComp1)
	if err == nil {
		t.Error("The record already exists")
	}

	_, err = em.Insert(newComp2)
	if err == nil {
		t.Error("Cif is unique, this cannot be added")
	}

	_, err = em.Insert(newComp3)
	if err == nil {
		t.Error("User mail does not exists, this insert must returns an error")
	}
}

func TestDB_Insert_Adding_Relations(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUC := &testusercompany{
		Email: "Plm",
		Cif:   "asd2",
	}

	_, err = em.Insert(newUC)
	if err != nil {
		t.Error(err)
	}

	newUC = &testusercompany{
		Email: "Plm1",
		Cif:   "asd4",
	}

	_, err = em.Insert(newUC)
	if err == nil {
		t.Error("Cif does not exist")
	}

	newUC = &testusercompany{
		Email: "Plm1123",
		Cif:   "asd1",
	}

	_, err = em.Insert(newUC)
	if err == nil {
		t.Error("User does not exist")
	}
}

func TestDB_First(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &testuser{
		Email: "Plm",
	}

	em.First(newUser, "", nil)

	if newUser.DNI != 123 {
		t.Error("DNI Unmatch")
	}
}

func TestDB_First_Not_Found(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &testuser{
		Email: "Plm245",
	}

	err = em.First(newUser, "", nil)
	if err == nil {
		t.Log(newUser)
		t.Error("This user does not exist")
	}

	err = em.First(newUser, "Admin = 0", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Find(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)
	foundUsers := make([]testuser, 0)

	err = em.Find(&foundUsers, "", nil)
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 3 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}

	where := "Admin = 1"

	foundUsers = make([]testuser, 0)

	err = em.Find(&foundUsers, where, nil)
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 2 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}

}

func TestDB_Find_Not_Found(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	foundUsers := make([]testuser, 0)

	err = em.Find(&foundUsers, "Admin = 3", nil)
	if err == nil {
		t.Error("Find must return an error")
	}
}

func TestDB_Remove(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &testuser{
		Email: "Plm2",
	}

	rs, err := em.Remove(newUser, "", nil)
	if err != nil {
		t.Error(err)
	}

	if rs.NumRecordsAffected != 1 {
		t.Error("Error removing existing record")
	}
}

func TestDB_Remove_Not_Found(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &testuser{
		Email: "Plm2421233",
	}

	rs, _ := em.Remove(newUser, "", nil)

	if rs.NumRecordsAffected != 0 {
		t.Error("Remove must returns an error because the record does not exist")
	}
}

func TestDB_Remove_Relation(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUC := &testusercompany{
		Email: "Plm",
		Cif:   "asd2",
	}

	_, err = em.Remove(newUC, "", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_DropTable(t *testing.T) {
	em, err := GetEntityManager(persistenceUnitItTest)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.DropTable(&testusercompany{})
	if err != nil {
		t.Error(err)
	}
	err = em.DropTable(&testuser{})
	if err != nil {
		t.Error(err)
	}
	err = em.DropTable(&testcompany{})
	if err != nil {
		t.Error(err)
	}

	_, err = em.Model(&testuser{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = em.Model(&testcompany{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = em.Model(&testusercompany{})
	if err == nil {
		t.Error("Model still exists")
	}
}
