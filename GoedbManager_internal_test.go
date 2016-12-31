package goedb

import (
	"testing"
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

var db *DB

func TestOpen(t *testing.T) {
	db = NewGoeDB()
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()
}

func TestDB_Migrate(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	err = db.Migrate(&TestUser{})
	if err != nil {
		t.Error(err)
	}

	if _, ok := db.tables["TestUser"]; !ok {
		t.Log(db.tables)
		t.Error("Migrate storage failed")
	}

	if db.tables["TestUser"].name == "" {
		t.Log(db.tables)
		t.Error("Table name unvalid")
	}

	if db.tables["TestUser"].columns == nil{
		t.Log(db.tables)
		t.Error("Migrate columns failed")
	}

	for key, value := range db.tables["TestUser"].columns {
		switch key{
		case 0:
			if !(value.title == "Email" && value.ctype == "string" && value.pk){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.title == "Password" && value.ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.title == "Role" && value.ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 3:
			if !(value.title == "DNI" && value.ctype == "int" && value.unique){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}

	err = db.Migrate(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	for key, value := range db.tables["TestCompany"].columns {
		switch key{
		case 0:
			if !(value.title == "UserEmail" && value.ctype == "string" && value.fk){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.title == "Name" && value.ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.title == "Cif" && value.ctype == "string" && value.pk){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}


	err = db.Migrate(&TestUserCompany{})
	if err != nil {
		t.Error(err)
	}

	for key, value := range db.tables["TestUserCompany"].columns {
		switch key{
		case 0:
			if !(value.title == "Email" && value.ctype == "string" && value.pk && value.fk){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.title == "Cif" && value.ctype == "string" && value.pk && value.fk){
				t.Log(value)
				t.Error("Column not valid")
			}
		}
	}
}

func TestDB_Model(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	user, err := db.Model(&TestUser{})
	if user.name != "TestUser" || len(user.columns) == 0{
		t.Error("Error getting db model")
	}

	company, err := db.Model(&TestCompany{})
	if company.name != "TestCompany" || len(company.columns) == 0{
		t.Error("Error getting db model")
	}
}

func TestDB_Model_Not_Found(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	_, err = db.Model(&OtherStruct{})
	if err == nil {
		t.Error("The result must has a error because the struct was not created")
	}
}


func TestDB_Insert(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

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

	_, err = db.Insert(newUser1)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Insert(newUser2)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Insert(newUser3)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Insert_with_FKs(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

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

	_, err = db.Insert(newComp1)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Insert(newComp2)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Insert_Constraints(t *testing.T) {

	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

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

	_, err = db.Insert(newComp1)
	if err == nil {
		t.Error("The record already exists")
	}

	_, err = db.Insert(newComp2)
	if err == nil {
		t.Error("Cif is unique, this cannot be added")
	}

	_, err = db.Insert(newComp3)
	if err == nil {
		t.Error("User mail does not exists, this insert must returns an error")
	}
}

func TestDB_Insert_Adding_Relations(t *testing.T) {

	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUC := &TestUserCompany{
		Email:"Plm",
		Cif:"asd2",
	}

	_, err = db.Insert(newUC)
	if err != nil {
		t.Error(err)
	}

	newUC = &TestUserCompany{
		Email:"Plm1",
		Cif:"asd4",
	}

	_, err = db.Insert(newUC)
	if err == nil {
		t.Error("Cif does not exist")
	}

	newUC = &TestUserCompany{
		Email:"Plm1123",
		Cif:"asd1",
	}

	_, err = db.Insert(newUC)
	if err == nil {
		t.Error("User does not exist")
	}
}

func TestDB_First(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := &TestUser{
		Email:"Plm",
	}

	db.First(newUser, "")

	if newUser.DNI != 123 {
		t.Error("DNI Unmatch")
	}
}

func TestDB_First_Not_Found(t *testing.T){
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := &TestUser{
		Email:"Plm245",
	}

	err = db.First(newUser, "")
	if err == nil {
		t.Log(newUser)
		t.Error("This user does not exist")
	}

	err = db.First(newUser, "Admin = 0")
	if err != nil {
		t.Error(err)
	}
}

func TestDB_Find(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	foundUsers := make([]TestUser, 0)


	err = db.Find(&foundUsers, "")
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 3 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}


	where := "Admin = 1"

	foundUsers = make([]TestUser, 0)

	err = db.Find(&foundUsers, where)
	if err != nil {
		t.Error(err)
	}

	if len(foundUsers) != 2 {
		t.Log(foundUsers)
		t.Error("Find not working")
	}

}

func TestDB_Find_Not_Found(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	foundUsers := make([]TestUser, 0)

	err = db.Find(&foundUsers, "Admin = 3")
	if err == nil {
		t.Error("Find must return an error")
	}
}

func TestDB_Remove(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := &TestUser{
		Email:"Plm2",
	}

	rs, err := db.Remove(newUser)
	if err != nil {
		t.Error(err)
	}

	if count, _ := rs.RowsAffected(); count != 1 {
		t.Error("Error removing existing record")
	}
}

func TestDB_Remove_Not_Found(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := &TestUser{
		Email:"Plm2421233",
	}

	rs, err := db.Remove(newUser)

	if count, _ := rs.RowsAffected(); count != 0 {
		t.Error("Remove must returns an error because the record does not exist")
	}
}

func TestDB_Remove_Relation(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUC := &TestUserCompany{
		Email:"Plm",
		Cif:"asd2",
	}

	_, err = db.Remove(newUC)
	if err != nil {
		t.Error(err)
	}
}

func TestDB_DropTable(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	db.DropTable(&TestUser{})
	db.DropTable(&TestCompany{})
	db.DropTable(&TestUserCompany{})
}
