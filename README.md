# GOEDB (Go Easy Database Manager)

Goedb is a simple ORM.

# How To Use

### Describing persistence.json

The file persistence.json will be used to define the datasource used at certain moment. It must be defined in the same directory as your program. A example of persistence.json is as follows:
```
    {
      "datasources":[
        {
          "name": "testSQLite3",
          "driver": "sqlite3",
          "url": "./test.db"
        }
      ]
    }
```

Currently multiple datasources can be defined. The name will be used as index to get the entity manager instance.

### Using Goedb

Once datasource is defined the next step is to get an instance of a entity manager. It requires the name of the datasource as input.

```
    em, err := GetEntityManager("testSQLite3")
	if err != nil {
		t.Error(err)
	}
```

In the current version, entity manager functionality is as follows:

```
type EntityManager interface {
	Open(driver string, params string) error
	Close() error
	Migrate(i interface{}) error
	DropTable(i interface{}) error
	Model(i interface{})(GoedbTable, error)
	Insert(i interface{}) (GoedbResult, error)
	Remove(i interface{}) (GoedbResult, error)
	First(i interface{}, params string) error
	Find(i interface{}, params string) error
}
```



# License

Copyright 2017 plopezm <pabloplm@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
