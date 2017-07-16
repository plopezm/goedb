package manager

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
)

type TestUser struct{
	Email		string	`goedb:"pk"`
	Password	string
	Role		string
	DNI		int	`goedb:"unique"`
	Admin		bool
}

type TestCompany struct {
	UserEmail	string	`goedb:"fk=TestUser(Email)"`
	Name		string
	Cif		string 	`goedb:"pk"`
}

type TestUserCompany struct {
	Email 		string 	`goedb:"pk,fk=TestUser(Email)"`
	Cif 		string	`goedb:"pk,fk=TestCompany(Cif)"`
}

type OtherStruct struct {
	Asd 	string
	Other	string
}

var dbSqlDriver *GoedbSQLDriver

func TestGoedbSQLDriver_Open(t *testing.T) {
	dbSqlDriver = new(GoedbSQLDriver)

	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error(err)
	}
	defer dbSqlDriver.Close()
}

func TestGoedbSQLDriver_Migrate(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error(err)
	}
	defer dbSqlDriver.Close()

	err = dbSqlDriver.Migrate(&TestUser{})
	if err != nil {
		t.Error(err)
	}

	if _, ok := dbSqlDriver.tables["TestUser"]; !ok {
		t.Log(dbSqlDriver.tables)
		t.Error("Migrate storage failed")
	}

	if dbSqlDriver.tables["TestUser"].Name == "" {
		t.Log(dbSqlDriver.tables)
		t.Error("Table name unvalid")
	}

	if dbSqlDriver.tables["TestUser"].Columns == nil{
		t.Log(dbSqlDriver.tables)
		t.Error("Migrate columns failed")
	}

	for key, value := range dbSqlDriver.tables["TestUser"].Columns {
		switch key{
		case 0:
			if !(value.Title == "Email" && value.ColumnType == "string" && value.PrimaryKey){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Title == "Password" && value.ColumnType == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.Title == "Role" && value.ColumnType == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 3:
			if !(value.Title == "DNI" && value.ColumnType == "int" && value.Unique){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

	err = dbSqlDriver.Migrate(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	for key, value := range dbSqlDriver.tables["TestCompany"].Columns {
		switch key{
		case 0:
			if !(value.Title == "UserEmail" && value.ColumnType == "string" && value.ForeignKey){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Title == "Name" && value.ColumnType == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.Title == "Cif" && value.ColumnType == "string" && value.PrimaryKey){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}


	err = dbSqlDriver.Migrate(&TestUserCompany{})
	if err != nil {
		t.Error(err)
	}

	for key, value := range dbSqlDriver.tables["TestUserCompany"].Columns {
		switch key{
		case 0:
			if !(value.Title == "Email" && value.ColumnType == "string" && value.PrimaryKey && value.ForeignKey){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.Title == "Cif" && value.ColumnType == "string" && value.PrimaryKey && value.ForeignKey){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}
}

func TestGoedbSQLDriver_Model(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	user, err := dbSqlDriver.Model(&TestUser{})
	if user.Name != "TestUser" || len(user.Columns) == 0{
		t.Error("Error getting db model")
	}

	company, err := dbSqlDriver.Model(&TestCompany{})
	if company.Name != "TestCompany" || len(company.Columns) == 0{
		t.Error("Error getting db model")
	}
}

func TestGoedbSQLDriver_Model_Not_Found(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	_, err = dbSqlDriver.Model(&OtherStruct{})
	if err == nil {
		t.Error("The result must has a error because the struct was not created")
	}
}


func TestGoedbSQLDriver_Insert(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUser1 := &TestUser{
		Email:"Plm",
		Password:"asd",
		Role: "asd",
		DNI:123,
		Admin: true,
	}

	newUser2 := &TestUser{
		Email:"Plm2",
		Password:"asd",
		Role: "asd",
		DNI:1234,
		Admin: true,
	}

	newUser3 := &TestUser{
		Email:"Plm3",
		Password:"asd",
		Role: "asd",
		DNI:1235,
		Admin: false,
	}

	_, err = dbSqlDriver.Insert(newUser1)
	if err != nil {
		t.Error(err)
	}

	_, err = dbSqlDriver.Insert(newUser2)
	if err != nil {
		t.Error(err)
	}

	_, err = dbSqlDriver.Insert(newUser3)
	if err != nil {
		t.Error(err)
	}
}

func TestGoedbSQLDriver_Insert_with_FKs(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newComp1 := &TestCompany{
		UserEmail:"Plm",
		Name:"asd",
		Cif: "asd1",
	}

	newComp2 := &TestCompany{
		UserEmail:"Plm",
		Name:"asd",
		Cif: "asd2",
	}

	_, err = dbSqlDriver.Insert(newComp1)
	if err != nil {
		t.Error(err)
	}

	_, err = dbSqlDriver.Insert(newComp2)
	if err != nil {
		t.Error(err)
	}
}

func TestGoedbSQLDriver_Insert_Constraints(t *testing.T) {

	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newComp1 := &TestCompany{
		UserEmail:"Plm",
		Name:"asd",
		Cif: "asd1",
	}

	newComp2 := &TestCompany{
		UserEmail:"Plm",
		Name:"asd123",
		Cif: "asd2",
	}

	newComp3 := &TestCompany{
		UserEmail:"Plm23455",
		Name:"asd",
		Cif: "asd4",
	}

	_, err = dbSqlDriver.Insert(newComp1)
	if err == nil {
		t.Error("The record already exists")
	}

	_, err = dbSqlDriver.Insert(newComp2)
	if err == nil {
		t.Error("Cif is unique, this cannot be added")
	}

	_, err = dbSqlDriver.Insert(newComp3)
	if err == nil {
		t.Error("User mail does not exists, this insert must returns an error")
	}
}

func TestGoedbSQLDriver_Insert_Adding_Relations(t *testing.T) {

	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUC := &TestUserCompany{
		Email:"Plm",
		Cif:"asd2",
	}

	_, err = dbSqlDriver.Insert(newUC)
	if err != nil {
		t.Error(err)
	}

	newUC = &TestUserCompany{
		Email:"Plm1",
		Cif:"asd4",
	}

	_, err = dbSqlDriver.Insert(newUC)
	if err == nil {
		t.Error("Cif does not exist")
	}

	newUC = &TestUserCompany{
		Email:"Plm1123",
		Cif:"asd1",
	}

	_, err = dbSqlDriver.Insert(newUC)
	if err == nil {
		t.Error("User does not exist")
	}
}

func TestGoedbSQLDriver_First(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUser := &TestUser{
		Email:"Plm",
	}

	dbSqlDriver.First(newUser, "")

	if newUser.DNI != 123 {
		t.Error("DNI Unmatch")
	}
}

func TestGoedbSQLDriver_First_Not_Found(t *testing.T){
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUser := &TestUser{
		Email:"Plm245",
	}

	err = dbSqlDriver.First(newUser, "")
	if err == nil {
		t.Log(newUser)
		t.Error("This user does not exist")
	}

	err = dbSqlDriver.First(newUser, "Admin = 0")
	if err != nil {
		t.Error(err)
	}
}

func TestGoedbSQLDriver_Find(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	foundUsers := make([]TestUser, 0)


	err = dbSqlDriver.Find(&foundUsers, "")
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 3 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}


	where := "Admin = 1"

	foundUsers = make([]TestUser, 0)

	err = dbSqlDriver.Find(&foundUsers, where)
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 2 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}

}

func TestGoedbSQLDriver_Find_Not_Found(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	foundUsers := make([]TestUser, 0)

	err = dbSqlDriver.Find(&foundUsers, "Admin = 3")
	if err == nil {
		t.Error("Find must return an error")
	}
}

func TestGoedbSQLDriver_Remove(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUser := &TestUser{
		Email:"Plm2",
	}

	rs, err := dbSqlDriver.Remove(newUser)
	if err != nil {
		t.Error(err)
	}

	if rs.NumRecordsAffected != 1 {
		t.Error("Error removing existing record")
	}
}

func TestGoedbSQLDriver_Remove_Not_Found(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUser := &TestUser{
		Email:"Plm2421233",
	}

	rs, err := dbSqlDriver.Remove(newUser)

	if rs.NumRecordsAffected != 0 {
		t.Error("Remove must returns an error because the record does not exist")
	}
}

func TestGoedbSQLDriver_Remove_Relation(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	newUC := &TestUserCompany{
		Email:"Plm",
		Cif:"asd2",
	}

	_, err = dbSqlDriver.Remove(newUC)
	if err != nil {
		t.Error(err)
	}
}

func TestGoedbSQLDriver_DropTable(t *testing.T) {
	err := dbSqlDriver.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer dbSqlDriver.Close()

	err = dbSqlDriver.DropTable(&TestUserCompany{})
	if err != nil {
		t.Error(err)
	}
	err = dbSqlDriver.DropTable(&TestUser{})
	if err != nil {
		t.Error(err)
	}
	err = dbSqlDriver.DropTable(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	_, err = dbSqlDriver.Model(&TestUser{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = dbSqlDriver.Model(&TestCompany{})
	if err == nil {
		t.Error("Model still exists")
	}

	_, err = dbSqlDriver.Model(&TestUserCompany{})
	if err == nil {
		t.Error("Model still exists")
	}
}
