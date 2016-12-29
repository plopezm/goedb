package goedb

import (
	"testing"
)

type TestUser struct{
	Email		string	`goedb:"pk"`
	Password	string
	Role		string
	DNI		int
	Admin		bool
}

type TestCompany struct {
	UserEmail	string	`goedb:"pk,fk=User(Email)"`
	Name		string
	Cif		string
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
		}
	}

	err = db.Migrate(&TestCompany{})
	if err != nil {
		t.Error(err)
	}

	for key, value := range db.tables["TestCompany"].columns {
		switch key{
		case 0:
			if !(value.title == "UserEmail" && value.ctype == "string" && value.pk && value.fk){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 1:
			if !(value.title == "Name" && value.ctype == "string"){
				t.Log(value)
				t.Error("Column not valid")
			}
		case 2:
			if !(value.title == "Cif" && value.ctype == "string"){
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

	user := db.Model(&TestUser{})
	if user.name != "TestUser" || len(user.columns) == 0{
		t.Error("Error getting db model")
	}

	company := db.Model(&TestCompany{})
	if company.name != "TestCompany" || len(company.columns) == 0{
		t.Error("Error getting db model")
	}
}

/*func TestMigrate(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := TestUser{
		Email:"Plm",
		Password:"asd",
		Role: "asd",
		DNI:123,
		Admin: true,
	}

	_, err = db.Model(&newUser).Insert(&newUser)
	if err != nil{
		t.Error(err)
	}
	//t.Log("Insert Result: ",result)

	var userFinded TestUser
	userFinded.Email ="Plm"
	err = db.Model(&TestUser{}).First(&userFinded, nil)
	if err != nil {
		t.Error(err)
	}

	if userFinded.DNI != newUser.DNI {
		t.Error("First: DNI does not match")
	}

	usersFound := make([]TestUser, 0)
	err = db.Model(&TestUser{}).Find(&usersFound, nil)
	if err != nil {
		t.Error(err)
	}

	if usersFound[0].DNI != newUser.DNI {
		t.Error("First: DNI does not match")
	}

	_, err = db.Model(&TestUser{}).Remove(&newUser)
	if err != nil{
		t.Error(err)
	}
	//t.Log("Delete Result: ",result)
}*/

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

func TestDB_Remove(t *testing.T) {
	err := db.Open("sqlite3", "./test.db")
	if err != nil{
		t.Error("DB couldn't be open")
	}
	defer db.Close()

	newUser := &TestUser{
		Email:"Plm2",
	}

	_, err = db.Remove(newUser)
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
}
