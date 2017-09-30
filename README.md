# GOEDB (Go Easy Database Manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/plopezm/goedb)](https://goreportcard.com/report/github.com/plopezm/goedb) [![Build Status](https://travis-ci.org/plopezm/goedb.svg)](https://travis-ci.org/plopezm/goedb) [![codecov](https://codecov.io/gh/plopezm/goedb/branch/master/graph/badge.svg)](https://codecov.io/gh/plopezm/goedb)

Goedb is a ORM for golang.


# How To Use

### Installation

This project uses [godep](https://github.com/golang/dep) for dependency management. To install the dependencies type the following:

```
    dep ensure
```

### Describing persistence.json

The file persistence.json will be used to define the datasource used at certain moment. It must be defined in the same directory as your program. A example of persistence.json is as follows:
```
    {
      "datasources":[
        {
          "name": "testSQLite3",
          "driver": "sqlite3",
          "url": "./test.db"
        },
        {
          "name": "testSQLite3Test",
          "driver": "sqlite3",
          "url": ":memory:"
        }
      ]
    }
```

Currently multiple datasources can be defined. The name will be used as index to get the entity manager instance.

### Using Goedb

Once datasource is defined the next step is to get an instance of a entity manager. It requires the name of the datasource as input.

```
	goedb.Initialize() // REQUIRED TO LOAD CONFIGURATION FROM persistence.json
	em, err = goedb.GetEntityManager("testSQLite3")
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
}
```

Now the manager is ready to work with him, for example:

```
	newUC := &TestUserCompany{
		Email:"Plm",
		Cif:"asd2",
	}

	_, err = em.Insert(newUC)
	if err != nil {
		t.Error(err)
	}
```

In the current version, entity manager functionality is as follows:

```
type EntityManager interface {
    SetSchema(schema string) (sql.Result, error)
    Open(driver string, params string, schema string) error
    Close() error
    Migrate(i interface{}) error
    DropTable(i interface{}) error
    Model(i interface{}) (metadata.GoedbTable, error)
    Insert(i interface{}) (GoedbResult, error)
    Update(i interface{}) (GoedbResult, error)
    Remove(i interface{}, where string, params map[string]interface{}) (GoedbResult, error)
    First(i interface{}, where string, params map[string]interface{}) error
    Find(i interface{}, where string, params map[string]interface{}) error
    TxBegin() (*sql.Tx, error)
}
```

# What is currently supported:

- For simple entities: All -> Tests in tests/Goedb_test.go
- For composed entities: All -> Tests in Goedb_ComplexStructs_test.go


- Polymorphism is not supported
- Named query supported in where clause. 
```
	err := em.First(soldier1, "TestSoldier.Name = :name", sql.Named("name", "Ryan"))
```

### Databases tested

- SQLite3
- PostgreSQL9

# License

Copyright 2017 plopezm <pabloplm@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
