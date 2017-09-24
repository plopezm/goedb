package tests

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/plopezm/goedb"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestUser struct {
	Email    string `goedb:"pk"`
	Password string
	Role     string
	DNI      int `goedb:"unique"`
	Admin    bool
}

type TestCompany struct {
	UserEmail string `goedb:"fk=TestUser(Email)"`
	Name      string
	Cif       string `goedb:"pk"`
}

type TestUserCompany struct {
	Email string `goedb:"pk,fk=TestUser(Email)"`
	Cif   string `goedb:"pk,fk=TestCompany(Cif)"`
}

type OtherStruct struct {
	Asd   string
	Other string
}

const PERSISTENCE_UNIT_IT_TEST = "testSQLite3"

func TestOpen(t *testing.T) {
	_, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Migrate(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.Migrate(&TestUser{})
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	err = em.Migrate(&TestUserCompany{})
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Model(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	user, _ := em.Model(&TestUser{})
	if user.Name != "TestUser" || len(user.Columns) == 0 {
		t.Error("Error getting db model")
	}

	company, _ := em.Model(&TestCompany{})
	if company.Name != "TestCompany" || len(company.Columns) == 0 {
		t.Error("Error getting db model")
	}
}

func TestDB_Model_Not_Found(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	_, err = em.Model(&OtherStruct{})
	if err == nil {
		t.Error("The result must has a error because the struct was not created")
	}
}

func TestDB_Insert(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser1 := &TestUser{
		Email:    "Plm",
		Password: "asd",
		Role:     "asd",
		DNI:      123,
		Admin:    true,
	}

	newUser2 := &TestUser{
		Email:    "Plm2",
		Password: "asd",
		Role:     "asd",
		DNI:      1234,
		Admin:    true,
	}

	newUser3 := &TestUser{
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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newComp1 := &TestCompany{
		UserEmail: "Plm",
		Name:      "asd",
		Cif:       "asd1",
	}

	newComp2 := &TestCompany{
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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newComp1 := &TestCompany{
		UserEmail: "Plm",
		Name:      "asd",
		Cif:       "asd1",
	}

	newComp2 := &TestCompany{
		UserEmail: "Plm",
		Name:      "asd123",
		Cif:       "asd2",
	}

	newComp3 := &TestCompany{
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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUC := &TestUserCompany{
		Email: "Plm",
		Cif:   "asd2",
	}

	_, err = em.Insert(newUC)
	if err != nil {
		t.Error(err)
	}

	newUC = &TestUserCompany{
		Email: "Plm1",
		Cif:   "asd4",
	}

	_, err = em.Insert(newUC)
	if err == nil {
		t.Error("Cif does not exist")
	}

	newUC = &TestUserCompany{
		Email: "Plm1123",
		Cif:   "asd1",
	}

	_, err = em.Insert(newUC)
	if err == nil {
		t.Error("User does not exist")
	}
}

func TestDB_First(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &TestUser{
		Email: "Plm",
	}

	em.First(newUser, "", nil)

	if newUser.DNI != 123 {
		t.Error("DNI Unmatch")
	}
}

func TestDB_First_Not_Found(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &TestUser{
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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)
	foundUsers := make([]TestUser, 0)

	err = em.Find(&foundUsers, "", nil)
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 3 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}

	where := "Admin = 1"

	foundUsers = make([]TestUser, 0)

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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	foundUsers := make([]TestUser, 0)

	err = em.Find(&foundUsers, "Admin = 3", nil)
	if err == nil {
		t.Error("Find must return an error")
	}
}

func TestDB_Remove(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &TestUser{
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
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUser := &TestUser{
		Email: "Plm2421233",
	}

	rs, _ := em.Remove(newUser, "", nil)

	if rs.NumRecordsAffected != 0 {
		t.Error("Remove must returns an error because the record does not exist")
	}
}

func TestDB_Remove_Relation(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	newUC := &TestUserCompany{
		Email: "Plm",
		Cif:   "asd2",
	}

	_, err = em.Remove(newUC, "", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_DropTable(t *testing.T) {
	em, err := goedb.GetEntityManager(PERSISTENCE_UNIT_IT_TEST)
	assert.Nil(t, err)
	assert.NotNil(t, em)

	err = em.DropTable(&TestUserCompany{})
	if err != nil {
		t.Error(err)
	}
	err = em.DropTable(&TestUser{})
	if err != nil {
		t.Error(err)
	}
	err = em.DropTable(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	_, err = em.Model(&TestUser{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = em.Model(&TestCompany{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = em.Model(&TestUserCompany{})
	if err == nil {
		t.Error("Model still exists")
	}
}
