package database



// Result is the result for some operation in database
type Result struct {
	NumRecordsAffected int64
	LastInsertId       int64
}
